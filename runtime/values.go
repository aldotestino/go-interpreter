package runtime

type ValueType string

const (
	NumberVT ValueType = "number"
)

type RuntimeValue interface {
	GetType() ValueType
	GetValue() any
}

// IntValue

type NumberValue struct {
	Type  ValueType
	Value float64
}

func NewNumberValue(num float64) *NumberValue {

	return &NumberValue{
		Type:  NumberVT,
		Value: num,
	}
}

func (nv *NumberValue) GetType() ValueType {
	return nv.Type
}

func (nv *NumberValue) GetValue() any {
	return nv.Value
}
