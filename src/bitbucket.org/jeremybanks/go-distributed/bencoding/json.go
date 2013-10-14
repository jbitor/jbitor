package bencoding

import (
	"bytes"
	"unicode/utf8"
)

// Converts Bencoded values to/from equivalent more-generic structures which
// can be encoded/decoded by encodings/json.
//
// Bencoded bytes strings are transformed into characters strings who
// codepoints correspond to the byte value.
//
// If a JSON Number contains a non-integer value, an error is returned.

func (dict Dict) ToJsonable() (jval interface{}, err error) {
	result := make(map[string]interface{})

	for key, value := range dict {
		jsonableKey, err := key.ToJsonable()
		if err != nil {
			return nil, err
		}

		jsonableKeyStr, ok := jsonableKey.(string)
		if !ok {
			panic("How isn't this a string?")
		}

		jsonableValue, err := value.ToJsonable()
		if err != nil {
			return nil, err
		}

		result[jsonableKeyStr] = jsonableValue
	}

	return result, nil
}

func (list List) ToJsonable() (jval interface{}, err error) {
	result := make([]interface{}, len(list))

	for i, value := range list {
		jsonableValue, err := value.ToJsonable()
		if err != nil {
			return nil, err
		}

		result[i] = jsonableValue
	}

	return result, nil
}

func (n Int) ToJsonable() (jval interface{}, err error) {
	return int64(n), nil
}

func (str String) ToJsonable() (jval interface{}, err error) {
	encodedStringPieces := make([][]byte, 0)

	for _, byteValue := range str {
		byteAsUtf8Char := make([]byte, 4)

		n := utf8.EncodeRune(byteAsUtf8Char, rune(byteValue))

		encodedStringPieces = append(encodedStringPieces, byteAsUtf8Char[:n])
	}

	return string(bytes.Join(encodedStringPieces, []byte(""))), nil
}

func FromJsonable(jval interface{}) (bval *Bencodable, err error) {
	logger.Fatalf("FromJsonable not implemented.\n")
	return
}
