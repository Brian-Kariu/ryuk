package flags

type ryukFlag struct {
	name       string
	value      string
	hasChanged bool
	valueType  string
}

type flagValue interface {
	Name() string
	ValueString() string
	ValueType() string
}

func (f ryukFlag) Name() string {
	return f.name
}

func (f ryukFlag) ValueString() string {
	return f.value
}

func (f ryukFlag) ValueType() string {
	return f.valueType
}

func NewFlag(name, value, valueType string) ryukFlag {
	return ryukFlag{name: name, value: value, valueType: valueType}
}
