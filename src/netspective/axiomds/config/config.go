package config

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"netspective/axiomds/rdbms"
)

type configDirectiveHandler func(c *Configuration, directive string, data interface {}) error

func configureDataSource(c *Configuration, directive string, data interface {}) error {
	structure, isMap := data.(map[string] interface{})
	if(! isMap) {
		return fmt.Errorf("Directive expects a hash, got %T instead", data)
	}

	name := c.getStringFromMap(structure, "", true, "Name")
	dataSource, err := rdbms.NewDataSource(name,
		c.getStringFromMap(structure, "", true, "Driver"),
		c.getStringFromMap(structure, "", true, "Connection Parameters", "Parameters"))

	if(err == nil) {
		c.dataSources[name] = *dataSource
	}
	return err
}

func configureSqlDml(c *Configuration, directive string, data interface {}) error {
	structure, isMap := data.(map[string] interface{})
	if(! isMap) {
		return fmt.Errorf("Directive expects a hash, got %T instead", data)
	}

	name := c.getStringFromMap(structure, "", true, "Name")
	sqlDml, err := rdbms.NewSqlDmlStatement(name, 
		c.getStringFromMap(structure, "", true, "Table Name"),
		c.getStringFromMap(structure, "Auto", false, "Command"),
		c.getStringListFromMap(structure, nil, false, "Columns", "Column", "Allowed Columns", "Allowed Column"),
		c.getStringListFromMap(structure, nil, false, "Primary Key Column", "Primary Key Columns"),
		c.getStringFromMap(structure, "Primary Keys", false, "Return on Success"))

	if(err == nil) {
		c.sqlDml[name] = *sqlDml
	}

	return err
}

func configureSqlDql(c *Configuration, directive string, data interface {}) error {
	structure, isMap := data.(map[string] interface{})
	if(! isMap) {
		return fmt.Errorf("Directive expects a hash, got %T instead", data)
	}

	name := c.getStringFromMap(structure, "", true, "Name")
	sqlDql, err := rdbms.NewSqlDqlStatement(name,
		c.getStringFromMap(structure, "", true, "DQL"),
		c.getStringListFromMap(structure, nil, false, "Parameters"),
		c.getStringFromMap(structure, "Multiple Rows as Hash List", false, "Return on Success"))

	if(err == nil) {
		c.sqlDql[name] = *sqlDql
	}

	return err
}

type include struct {
	path string
	isDir bool
}

type Configuration struct {
	// TODO: includes []include
	handlers map[string] configDirectiveHandler
	dataSources map[string] rdbms.DataSource
	sqlDml map[string] rdbms.SqlDmlStatement
	sqlDql map[string] rdbms.SqlDqlStatement

	activeFile string
	activeDirectiveIndex int
	activeDirectiveName string
}

func NewConfiguration() *Configuration {
	instance := new(Configuration)

	instance.handlers = make(map[string] configDirectiveHandler)
	instance.registerConfigHandler(configureDataSource, "RDBMS Data Source")
	instance.registerConfigHandler(configureSqlDml, "SQL DML")
	instance.registerConfigHandler(configureSqlDql, "SQL DQL")

	instance.dataSources = make(map[string] rdbms.DataSource)
	instance.sqlDml = make(map[string] rdbms.SqlDmlStatement)
	instance.sqlDql = make(map[string] rdbms.SqlDqlStatement)

	instance.activeFile = ""
	instance.activeDirectiveIndex = -1
	instance.activeDirectiveName = ""
	return instance
}

func (c Configuration) DataSources() map[string] rdbms.DataSource { return c.dataSources }

func (c *Configuration) startFile(file string) ([]byte, error) {
	c.activeFile = file
	content, err := ioutil.ReadFile(file)
	return content, err
}

func (c *Configuration) finishFile() {
	c.activeFile = ""
}

func (c *Configuration) startDirectiveInActiveFile(index int) {
	c.activeDirectiveIndex = index
}

func (c *Configuration) startDirectiveNameInActiveFile(name string) {
	c.activeDirectiveName = name
}

func (c *Configuration) finishDirectiveInActiveFile() {
	c.activeDirectiveIndex = -1
	c.activeDirectiveName = ""
}

func (c *Configuration) getStringFromMap(data map[string] interface{}, defaultValue string, required bool, keys ...string) string {
	for keyIndex := range keys {
		key := keys[keyIndex]
		value, ok := data[key]
		if(ok) {
			text, ok := value.(string)
			if(ok) {
				return text;
			}
		}
	}

	if(required) {
		c.logError(fmt.Sprintf("Expected a string value with one of these keys: %s", keys))
	}
	return defaultValue;
}

func (c *Configuration) getStringListFromMap(data map[string] interface{}, defaultValue []string, required bool, keys ...string) []string {
	for keyIndex := range keys {
		key := keys[keyIndex]
		value, ok := data[key]
		if(ok) {
			textList, ok := value.([]string)
			if(ok) {
				return textList;
			} else {
				text, ok := value.(string)
				if(ok) {
					return []string { text };
				}
			}
		}
	}

	if(required) {
		c.logError(fmt.Sprintf("Expected a string list value with one of these keys: %s", keys))
	}
	return defaultValue;
}

func (c *Configuration) registerConfigHandler(handler configDirectiveHandler, names ... string) {
	if(handler == nil) {
		panic("configDirectiveHandler is nil")
	}

	for i := range names {
		name := names[i]

		_, found := c.handlers[name]
		if(found) {
			panic(fmt.Sprintf("[registerConfigHandler] Duplicate configDirectiveHandler: %s", name))
		}

		c.handlers[name] = handler
	}
}

func (c *Configuration) logError(message interface{}) {
	fmt.Printf("[%s, %d, '%s'] %s\n", c.activeFile, c.activeDirectiveIndex, c.activeDirectiveName, message)
}

func (c *Configuration) Configure(file string) bool {
	content, err := c.startFile(file)
	if(err != nil) {
		c.logError(fmt.Sprintf("Error reading file: %s", err))
		return false;
	}

	var instance interface{}
	err = json.Unmarshal(content, &instance)
	if(err != nil) {
		c.logError(fmt.Sprintf("Error parsing file as JSON: %s", err))
		return false;
	}

	directivesBlock, directivesFound := instance.(map[string] interface{})
	if(! directivesFound) {
		c.logError(fmt.Sprintf("Unmarshalling JSON does not result in map[string] interface{}"))
		return false;
	}

	configureBlock, configureFound := directivesBlock["configure"];
	if(! configureFound) {
		c.logError(fmt.Sprintf("JSON does not contain a \"configure\" directive with an array of directives"))
		return false;
	}

	configureList, isConfigureArray := configureBlock.([] interface {});
	if(! isConfigureArray) {
		c.logError("\"configure\" value is not an array")
		return false;
	}

	for index := range configureList {
		c.startDirectiveInActiveFile(index)
		defer c.finishDirectiveInActiveFile()

		directive := configureList[index]
		_, isDirectiveMap := directive.(map[string] interface{})
		if(! isDirectiveMap) {
			c.logError("Directive is not a map[string] interface{}")
			return false;
		}

		for key, value := range directive.(map[string] interface{}) {
			c.startDirectiveNameInActiveFile(key);

			handler, handlerFound := c.handlers[key]
			if(! handlerFound) {
				c.logError("Directive does not have a handler")
				continue
			}

			err := handler(c, key, value)
			if(err != nil) {
				c.logError(err)
			}
		}
	}

	return true;
}
