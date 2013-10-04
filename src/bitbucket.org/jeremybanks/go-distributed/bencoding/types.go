package bencoding

type ValueType int

const (
	STRING ValueType = iota
	INTEGER
	LIST
	DICTIONARY
)

type Value struct {
	T     ValueType
	Value interface{}
}

type Bencodable interface {
	MarshalBencodingValue() (*Value, error)
}

type Marshaler interface {
	MarshalBencoding() ([]byte, error)
}

type Bdecodable interface {
	UnmarshalBencodingValue(*Value) error
}

type Unmarshaler interface {
	UnmarshalBencoding([]byte) error
}
