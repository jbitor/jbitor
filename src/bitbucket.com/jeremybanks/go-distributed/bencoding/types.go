package bencoding

type BValueType int

const (
	STRING BValueType = iota
	INTEGER
	LIST
	DICTIONARY
)

type BValue struct {
	t     BValueType
	value interface{}
}

type Bencodable interface {
	BValue() *BValue
}

// type BDecodable interface {
//  initFromBValue(*BValue) error
// }
