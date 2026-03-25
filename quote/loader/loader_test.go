package loader

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"testing"
	"testing/fstest"

	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer"
)

var inverseKeyTable = [256]byte{
	0x37, 0x6a, 0x09, 0x5e, 0x7a, 0xaf, 0xf5, 0xa4, 0xba, 0x78, 0x84, 0x58, 0x35, 0x1e, 0x6b, 0x0c,
	0x49, 0xc6, 0xc3, 0x44, 0x40, 0x9e, 0x6f, 0x65, 0xe4, 0xf6, 0xfe, 0x22, 0xe2, 0x95, 0xc7, 0x38,
	0xf0, 0x1a, 0x82, 0xe0, 0x5b, 0x2a, 0xd8, 0xe5, 0xce, 0x2f, 0x74, 0x25, 0xec, 0x59, 0xc0, 0x45,
	0x4b, 0x64, 0x43, 0xdc, 0xb0, 0xb9, 0x30, 0x6d, 0x28, 0xd1, 0x16, 0xbb, 0x66, 0x98, 0x92, 0x90,
	0x2c, 0xa7, 0xf1, 0x80, 0xc1, 0xd4, 0x8b, 0xd6, 0xdf, 0x24, 0x2d, 0xf7, 0xfb, 0x88, 0x4d, 0x3c,
	0x72, 0xf3, 0xdb, 0x2b, 0x93, 0x73, 0xef, 0x85, 0x83, 0xee, 0xc2, 0x8d, 0x5c, 0xb2, 0x0b, 0x94,
	0x3d, 0xa8, 0x3f, 0x1c, 0x4c, 0x6e, 0x03, 0x7b, 0x1d, 0x5a, 0x51, 0xa1, 0x70, 0x41, 0xd0, 0xaa,
	0xa0, 0x7e, 0xcd, 0xd5, 0x15, 0xa9, 0x18, 0x76, 0xc9, 0x7d, 0x7f, 0x0e, 0x3a, 0x99, 0xbf, 0xab,
	0x3b, 0x14, 0x3e, 0x9a, 0x04, 0xda, 0x02, 0xfd, 0x63, 0xd9, 0xfa, 0x9f, 0x4e, 0xe3, 0x61, 0xbe,
	0x07, 0x11, 0xa6, 0x1b, 0x19, 0x55, 0x8e, 0x77, 0x0a, 0x47, 0xe6, 0xf8, 0x0d, 0xcf, 0xd7, 0x33,
	0x23, 0x1f, 0xbc, 0x62, 0xde, 0x9b, 0x29, 0x53, 0x68, 0xe8, 0x21, 0xb6, 0x34, 0x52, 0x87, 0xcb,
	0x08, 0x79, 0xf4, 0x67, 0x69, 0x54, 0xe7, 0x86, 0xea, 0xb4, 0x20, 0x71, 0x01, 0xbd, 0x06, 0x31,
	0x00, 0x50, 0xc8, 0xb8, 0xac, 0x5d, 0x57, 0x7c, 0x89, 0xeb, 0xb7, 0x36, 0x8f, 0xf2, 0xe1, 0x56,
	0x81, 0x4a, 0xd2, 0x8c, 0xf9, 0xad, 0x60, 0xa5, 0x42, 0x10, 0x5f, 0x12, 0xb3, 0xff, 0x4f, 0xdd,
	0x46, 0x26, 0xa2, 0x17, 0xc5, 0x75, 0x91, 0x27, 0xb5, 0x8a, 0xd3, 0x13, 0x2e, 0xc4, 0xe9, 0x9d,
	0x97, 0x39, 0x32, 0x05, 0x0f, 0xca, 0xcc, 0x48, 0xfc, 0xae, 0x96, 0xed, 0x6c, 0x9c, 0xb1, 0xa3,
}

const (
	pass1XorA byte = 0x45
	pass1XorB byte = 0x71
	pass2XorA byte = 0x86
	pass2XorB byte = 0x23
)

var forwardKeyTable [256]byte

func init() {
	for i := 0; i < 256; i++ {
		forwardKeyTable[inverseKeyTable[i]] = byte(i)
	}
}

func xorEncode(data []byte, xorA, xorB byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = forwardKeyTable[b^xorB] ^ xorA
	}
	return out
}

func encodeTestPayload(plaintext []byte) []byte {
	pass1Encoded := xorEncode(plaintext, pass1XorA, pass1XorB)

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(pass1Encoded)
	w.Close()

	pass2Encoded := xorEncode(compressed.Bytes(), pass2XorA, pass2XorB)

	var header [16]byte
	copy(header[:4], "ONS2")
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(pass2Encoded)))
	binary.LittleEndian.PutUint32(header[8:12], uint32(len(plaintext)))
	binary.LittleEndian.PutUint32(header[12:16], 110)

	var out bytes.Buffer
	out.Write(header[:])
	out.Write(pass2Encoded)
	return out.Bytes()
}

func testParseFunc(quotes []dto.ParsedQuote, refs []lexer.SubtitleRef) ParseFunc {
	return func(lines []string) ([]dto.ParsedQuote, []lexer.SubtitleRef, []lexer.ValidationError) {
		return quotes, refs, nil
	}
}

func buildTestFS(path string, plaintext []byte) fstest.MapFS {
	return fstest.MapFS{
		path: &fstest.MapFile{Data: encodeTestPayload(plaintext)},
	}
}

func TestLoad_ReturnsQuotes(t *testing.T) {
	expected := []dto.ParsedQuote{
		{Text: "Hello", CharacterID: "10", Episode: 1},
		{Text: "World", CharacterID: "27", Episode: 1},
	}
	fs := buildTestFS("data/test.file", []byte("line1\nline2"))
	loader := New(fs, testParseFunc(expected, nil))

	result := loader.Load("en", "data/test.file")

	if len(result) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(result))
	}
	if result[0].Text != "Hello" {
		t.Errorf("quote 0 text: got %q, want %q", result[0].Text, "Hello")
	}
	if result[1].Text != "World" {
		t.Errorf("quote 1 text: got %q, want %q", result[1].Text, "World")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	fs := fstest.MapFS{}
	loader := New(fs, testParseFunc(nil, nil))

	result := loader.Load("en", "data/missing.file")

	if result != nil {
		t.Errorf("expected nil for missing file, got %d quotes", len(result))
	}
}

func TestLoad_InvalidEncodedData(t *testing.T) {
	fs := fstest.MapFS{
		"data/bad.file": &fstest.MapFile{Data: []byte("not valid ONS2 data here!")},
	}
	loader := New(fs, testParseFunc(nil, nil))

	result := loader.Load("en", "data/bad.file")

	if result != nil {
		t.Errorf("expected nil for invalid data, got %d quotes", len(result))
	}
}

func TestLoad_PassesDecodedLinesToParser(t *testing.T) {
	plaintext := []byte("first line\nsecond line\nthird line")
	fs := buildTestFS("data/test.file", plaintext)

	var capturedLines []string
	parse := func(lines []string) ([]dto.ParsedQuote, []lexer.SubtitleRef, []lexer.ValidationError) {
		capturedLines = lines
		return nil, nil, nil
	}
	loader := New(fs, parse)

	loader.Load("en", "data/test.file")

	if len(capturedLines) != 3 {
		t.Fatalf("expected 3 lines passed to parser, got %d", len(capturedLines))
	}
	if capturedLines[0] != "first line" {
		t.Errorf("line 0: got %q, want %q", capturedLines[0], "first line")
	}
	if capturedLines[2] != "third line" {
		t.Errorf("line 2: got %q, want %q", capturedLines[2], "third line")
	}
}

func TestLoad_ResolvesSubtitleRefs(t *testing.T) {
	assContent := "[Script Info]\nTitle: Test\n\n[V4+ Styles]\nFormat: Name\nStyle: Default\n\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\nDialogue: 0,0:00:00.00,0:00:05.00,Default,,0,0,0,,Welcome back.\nDialogue: 0,0:00:05.00,0:00:10.00,Default,,0,0,0,,Goodbye.\n"

	fs := fstest.MapFS{
		"data/test.file":      &fstest.MapFile{Data: encodeTestPayload([]byte("test"))},
		"data/sub/ending.ass": &fstest.MapFile{Data: []byte(assContent)},
	}

	refs := []lexer.SubtitleRef{
		{SubPath: `sub\ending.ass`, CharacterID: "00", AudioID: "end_test", Episode: 8},
	}
	loader := New(fs, testParseFunc(nil, refs))

	result := loader.Load("en", "data/test.file")

	if len(result) != 2 {
		t.Fatalf("expected 2 subtitle quotes, got %d", len(result))
	}
	if result[0].AudioID != "end_test_s0" {
		t.Errorf("quote 0 audioID: got %q, want %q", result[0].AudioID, "end_test_s0")
	}
	if result[1].AudioID != "end_test_s1" {
		t.Errorf("quote 1 audioID: got %q, want %q", result[1].AudioID, "end_test_s1")
	}
	if result[0].Episode != 8 {
		t.Errorf("quote 0 episode: got %d, want 8", result[0].Episode)
	}
}

func TestLoad_SubtitleRefMissingFile(t *testing.T) {
	fs := buildTestFS("data/test.file", []byte("test"))

	refs := []lexer.SubtitleRef{
		{SubPath: `sub\missing.ass`, CharacterID: "00", AudioID: "end_test", Episode: 8},
	}
	loader := New(fs, testParseFunc(nil, refs))

	result := loader.Load("en", "data/test.file")

	if len(result) != 0 {
		t.Errorf("expected 0 quotes when subtitle file missing, got %d", len(result))
	}
}

func TestLoad_CombinesParsedAndSubtitleQuotes(t *testing.T) {
	assContent := "[Script Info]\nTitle: Test\n\n[V4+ Styles]\nFormat: Name\nStyle: Default\n\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\nDialogue: 0,0:00:00.00,0:00:05.00,Default,,0,0,0,,Sub line.\n"

	fs := fstest.MapFS{
		"data/test.file":      &fstest.MapFile{Data: encodeTestPayload([]byte("test"))},
		"data/sub/ending.ass": &fstest.MapFile{Data: []byte(assContent)},
	}

	parsed := []dto.ParsedQuote{
		{Text: "Parsed quote", CharacterID: "10", Episode: 1},
	}
	refs := []lexer.SubtitleRef{
		{SubPath: `sub\ending.ass`, CharacterID: "00", AudioID: "end_test", Episode: 8},
	}
	loader := New(fs, testParseFunc(parsed, refs))

	result := loader.Load("en", "data/test.file")

	if len(result) != 2 {
		t.Fatalf("expected 2 total quotes (1 parsed + 1 subtitle), got %d", len(result))
	}
	if result[0].Text != "Parsed quote" {
		t.Errorf("quote 0: got %q, want %q", result[0].Text, "Parsed quote")
	}
	if result[1].Text != "Sub line." {
		t.Errorf("quote 1: got %q, want %q", result[1].Text, "Sub line.")
	}
}
