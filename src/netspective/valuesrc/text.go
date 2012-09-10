package valuesrc

type TextLiteral struct {
	literal string
}

func NewTextLiteral(literal string) (ValueSource, error) {
	instance := new(TextLiteral)
	instance.literal = literal
	return instance, nil
}

func (tl TextLiteral) EvaluateValue(vc ValueContext) Value {
	return tl;
}

func (tl TextLiteral) GetValue() interface{} {
	return tl.literal
}

func (tl TextLiteral) HasValue() bool {
	return true
}
