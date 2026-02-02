package encoder

import "testing"

func TestEncode(t *testing.T) {
	
	tests := []struct {
	    input    uint64
	    expected string
	}{
	    {0, "0"},
	    {61, "z"},         // 61 is the last char: lowercase 'z'
	    {62, "10"},
	    {12345, "3D7"},    // Note the capital D
	    {18446744073709551615, "LygHa16AHYF"}, 
	}

	for _, tc := range tests {
		got := Encode(tc.input)
		if got != tc.expected {
			t.Errorf("Encode(%d) failed: got %s, want %s", tc.input, got, tc.expected)
		}
	}
}