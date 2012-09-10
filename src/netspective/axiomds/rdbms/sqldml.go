package rdbms

import "fmt"

type SqlDmlCommand int
const (
	SqlDmlCommand_Auto SqlDmlCommand = iota
	SqlDmlCommand_Insert
	SqlDmlCommand_Update
	SqlDmlCommand_Delete
)
var SqlDmlCommandNames = map[string] SqlDmlCommand {
	"Auto" : SqlDmlCommand_Auto,
	"Insert" : SqlDmlCommand_Insert,
	"Update" : SqlDmlCommand_Update,
	"Delete" : SqlDmlCommand_Delete,
}

type SqlDmlReturnOnSuccess int
const (
	SqlDmlReturn_None SqlDmlReturnOnSuccess = iota
	SqlDmlReturn_PrimaryKeys
	SqlDmlReturn_PrimaryKeysAndParametersPassed
	SqlDmlReturn_ParametersPassed
	SqlDmlReturn_SelectColumnsFromTable
	SqlDmlReturn_SelectAllFromTable
	SqlDmlReturn_CustomSqlDql
)
var SqlDmlReturnOnSuccessNames = map[string] SqlDmlReturnOnSuccess {
	"None" : SqlDmlReturn_None,
	"Primary Keys" : SqlDmlReturn_PrimaryKeys,
	"Primary Key" : SqlDmlReturn_PrimaryKeys,
	"Primary Keys and Parameters Passed" : SqlDmlReturn_PrimaryKeysAndParametersPassed,
	"Primary Key and Parameters Passed" : SqlDmlReturn_PrimaryKeysAndParametersPassed,
	"Parameters Passed" : SqlDmlReturn_ParametersPassed,
	"Select Specific Columns from Table" : SqlDmlReturn_SelectColumnsFromTable,
	"Select All from Table Name" : SqlDmlReturn_SelectAllFromTable,
	//TODO: "SQL DQL" : SqlDmlReturn_CustomSqlDql,
}

type SqlDmlStatement struct {
	identity string
	tableName string
	command SqlDmlCommand
	allowedColumns []string
	primaryKeyColumns []string
	returnOnSuccess SqlDmlReturnOnSuccess
	//TODO: returnOnSuccessCustomDql SqlDqlStatement
}

func NewSqlDmlStatement(identity string, tableName string, commandText string, allowedColumns []string,
						primaryKeyColumns []string, returnOnSuccessText string) (*SqlDmlStatement, error) {
	instance := new(SqlDmlStatement)
	instance.identity = identity
	instance.tableName = tableName
	instance.allowedColumns = allowedColumns
	instance.primaryKeyColumns = primaryKeyColumns

	var err error = nil

	command, foundCmd := SqlDmlCommandNames[commandText]
	if(foundCmd) {
		instance.command = command
	} else {
		err = fmt.Errorf("commandText '%s' is not valid, expected one of: %s", commandText, SqlDmlCommandNames)
	}

	returnOnSuccess, foundROS := SqlDmlReturnOnSuccessNames[returnOnSuccessText]
	if(foundROS) {
		instance.returnOnSuccess = returnOnSuccess
	} else {
		err = fmt.Errorf("returnOnSuccessText '%s' is not valid, expected one of: %s", returnOnSuccessText, SqlDmlReturnOnSuccessNames)
	}

	return instance, err
}
