package bencoding

import "testing"

// Test cases

func TestIntegerEncoding(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"i3e":   int64(3),
		"i-3e":  int64(-3),
		"i6e":   int64(6),
		"i0e":   int64(0),
		"i16e":  Value{T: INTEGER, Value: int64(16)},  // already a Value
		"i17e":  &Value{T: INTEGER, Value: int64(17)}, // already a *Value
		"i-99e": bencodableInt32(-99),
	})

	testUnencodables(t, []interface{}{
		int32(99),  // wrong integer type
		uint64(99), // wrong integer type
	})
}

func TestListEncoding(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"l4:spam4:eggse": []interface{}{"spam", "eggs"},
	})
}

func TestDictionaryEncoding(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"d3:cow3:moo4:spam4:eggse": map[string]interface{}{"cow": "moo", "spam": "eggs"},
		"d4:spaml1:a1:bee":         map[string]interface{}{"spam": []interface{}{"a", "b"}},
	})
}

func TestEncodeUnverifiedExample(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"d6:lengthi512e4:miscd5:hello6:World!e4:name9:Test Data12:piece lengthi1024e6:pieces20:\x00234567890123456789\xFFe": map[string]interface{}{
			"piece length": int64(1024),
			"pieces":       "\x00234567890123456789\xFF",
			"name":         "Test Data",
			"length":       int64(512),
			"misc": map[string]interface{}{
				"hello": "World!",
			},
		},
	})
}
