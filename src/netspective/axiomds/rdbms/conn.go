package rdbms

type DataSource struct {
	name string
	driver string
	connParams string
}

func NewDataSource(name string, driver string, connParams string) (*DataSource, error) {
	instance := new(DataSource)
	instance.name = name
	instance.driver = driver
	instance.connParams = connParams
	return instance, nil
}

func (ds DataSource) Name() string { return ds.name }
func (ds DataSource) Driver() string { return ds.driver }
func (ds DataSource) ConnectionParameters() string { return ds.connParams }

