package rdbms

import "fmt"

type SqlDqlReturnOnSuccess int
const (
	SqlDqlReturn_SingleValue SqlDqlReturnOnSuccess = iota
	SqlDqlReturn_ListOfSingleValues
	SqlDqlReturn_SingleRowAsHash
	SqlDqlReturn_MultipleRowsAsMatrix
	SqlDqlReturn_MultipleRowsAsHashList
)
var SqlDqlReturnOnSuccessNames = map[string] SqlDqlReturnOnSuccess {
	"Single Value" : SqlDqlReturn_SingleValue,
	"List of Single Values" : SqlDqlReturn_ListOfSingleValues,
	"Single Row as Hash" : SqlDqlReturn_SingleRowAsHash,
	"Multiple Rows as Matrix" : SqlDqlReturn_MultipleRowsAsMatrix,
	"Multiple Rows as Hash List" : SqlDqlReturn_MultipleRowsAsHashList,
}

type SqlDql struct {
	parameters SqlStatementParameters
	// TODO: isTemplate bool
	dql string
	returnOnSuccess SqlDqlReturnOnSuccess
	dataSource *string
}

type SqlDqlStatement struct {
	identity string
	dqls []SqlDql
}

func (stmt SqlDqlStatement) IsConnectionSpecific() bool {
	return len(stmt.dqls) > 1
}

func NewSqlDql(dql string, parameterSpecs []string, returnOnSuccessText string, dataSource *string) (*SqlDql, error) {
	instance := new(SqlDql)
	instance.dql = dql
	instance.dataSource = dataSource

	params, paramsErr := NewSqlStatementParametersSimple(parameterSpecs);
	if(paramsErr != nil) {
		return instance, paramsErr
	}
	instance.parameters = *params

	returnOnSuccess, foundROS := SqlDqlReturnOnSuccessNames[returnOnSuccessText]
	if(foundROS) {
		instance.returnOnSuccess = returnOnSuccess
	} else {
		return instance, fmt.Errorf("returnOnSuccessText '%s' is not valid, expected one of: %s", returnOnSuccessText, SqlDqlReturnOnSuccessNames)
	}

	return instance, nil
}

func NewSqlDqlStatement(identity string, dqlText string, parameterSpecs []string,
                        returnOnSuccessText string) (*SqlDqlStatement, error) {
	instance := new(SqlDqlStatement)
	instance.identity = identity

	dql, err := NewSqlDql(dqlText, parameterSpecs, returnOnSuccessText, nil)
	if(err == nil) {
		return instance, err
	}
	instance.dqls = append(instance.dqls, *dql)
	return instance, nil
}
