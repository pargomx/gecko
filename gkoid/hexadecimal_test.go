package gkoid

import (
	"testing"
)

func TestHex_String(t *testing.T) {
	tests := []struct {
		id     Hex
		expect string
	}{
		{0x0, "0"},
		{0x1, "1"},
		{0xA, "a"},
		{0x10, "10"},
		{0x123, "123"},
		{0xDEADBEEF, "deadbeef"},
		{0xFFFFFFFFFFFFFFFF, "ffffffffffffffff"},
		{0x00000000000000FF, "ff"},
		{0x0000000000000100, "100"},
	}
	for _, tt := range tests {
		got := tt.id.String()
		if got != tt.expect {
			t.Errorf("Hex(%x).String() = %q, want %q", uint64(tt.id), got, tt.expect)
		}
	}
}

func BenchmarkHex_String(b *testing.B) {
	ids := []Hex{
		0x0, 0x1, 0xA, 0x10, 0x123, 0xDEADBEEF, 0xFFFFFFFFFFFFFFFF,
	}
	for i := 0; b.Loop(); i++ {
		_ = ids[i%len(ids)].String()
	}
}

// func BenchmarkHex_StringFmt(b *testing.B) {
// 	ids := []Hex{
// 		0x0, 0x1, 0xA, 0x10, 0x123, 0xDEADBEEF, 0xFFFFFFFFFFFFFFFF,
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = ids[i%len(ids)].StringFmt()
// 	}
// }

func TestParseHex(t *testing.T) {
	tests := []struct {
		input    string
		expected Hex
		wantErr  bool
	}{
		{"0", 0x0, false},
		{"1", 0x1, false},
		{"a", 0xa, false},
		{"10", 0x10, false},
		{"123", 0x123, false},
		{"deadbeef", 0xDEADBEEF, false},
		{"ffffffffffffffff", 0xFFFFFFFFFFFFFFFF, false},
		{"ff", 0xff, false},
		{"100", 0x100, false},
		{"", 0, true},
		{"zz", 0, true},
		{"fffffffffffffffff", 0, true}, // too long
	}
	for _, tt := range tests {
		got, err := ParseHex(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseHex(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if err == nil && got != tt.expected {
			t.Errorf("ParseHex(%q) = %x, want %x", tt.input, got, tt.expected)
		}
	}
}

func BenchmarkParseHex(b *testing.B) {
	inputs := []string{
		"0", "1", "a", "10", "123", "deadbeef", "ffffffffffffffff", "ff", "100",
	}
	for i := 0; b.Loop(); i++ {
		_, _ = ParseHex(inputs[i%len(inputs)])
	}
}

func TestNewHex(t *testing.T) {
	tests := []struct {
		digitos int
		wantErr bool
	}{
		{1, true},
		{2, false},
		{8, false},
		{16, false},
		{17, true},
		{0, true},
	}
	for _, tt := range tests {
		id, err := NewHex(tt.digitos)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewHex(%d) error = %v, wantErr %v", tt.digitos, err, tt.wantErr)
		}
		if err == nil {
			// Check that the string representation has at most hex digits
			s := id.String()
			if len(s) > tt.digitos {
				t.Errorf("NewHex(%d) generated %q with %d digits", tt.digitos, s, len(s))
			}
			// Check that it can be parsed
			parsed, err := ParseHex(s)
			if err != nil {
				t.Errorf("NewHex(%d) generated %q which failed to parse: %v", tt.digitos, s, err)
			} else if parsed != id {
				t.Errorf("NewHex(%d) generated %q which parsed to %x, want %x", tt.digitos, s, parsed, id)
			}
		}
	}
}

func BenchmarkNewHex(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		_, _ = NewHex(8)
	}
}
