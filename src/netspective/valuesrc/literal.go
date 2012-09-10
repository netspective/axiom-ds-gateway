package valuesrc

type Literal struct {
	literal interface{}
}

func NewLiteral(literal string) (ValueSource, error) {
	instance := new(Literal)
	instance.literal = literal
	return *instance, nil
}

func (l Literal) EvaluateValue(vc ValueContext) Value {
	return l;
}

func (l Literal) GetValue() interface{} {
	return l.literal
}

func (l Literal) HasValue() bool {
	return true
}
