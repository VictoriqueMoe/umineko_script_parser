# umineko_script_parser

A Go library for parsing Umineko no Naku Koro ni (When the Seagulls Cry) game script files. Extracts structured dialogue quotes with character attribution, episode metadata, voice audio references, sound effect associations, red/blue/gold/purple truth detection, and both plain text and HTML output.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Source Text                                    │
│  d [lv 0*"27"*"10100001"]`"{p:1:Without love, it cannot be seen.}"`[\]     │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           LEXER (lexer/)                                    │
│  Tokenises input into a stream of typed tokens                              │
│  • TokenCommand: "d"                                                        │
│  • TokenInlineCommand: "lv 0*\"27\"*\"10100001\""                           │
│  • TokenFormatTag: "p:1:Without love, it cannot be seen."                   │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          PARSER (lexer/)                                    │
│  Builds Abstract Syntax Tree from tokens                                    │
│                                                                             │
│  Script                                                                     │
│   └── Lines[]                                                               │
│        ├── EpisodeMarkerLine { Episode: 1, Type: "episode" }                │
│        ├── PresetDefineLine { ID: 1, Colour: "#FF0000" }                    │
│        └── DialogueLine                                                     │
│             ├── Command: "d"                                                │
│             └── Content[]                                                   │
│                  ├── VoiceCommand { CharacterID: "27", AudioID: "..." }     │
│                  └── FormatTag { Name: "p", Param: "1", Content: [...] }    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       VALIDATOR (lexer/)                                     │
│  Post-parse AST validation (non-fatal)                                      │
│                                                                             │
│  • Unknown format tags        • Missing voice command fields                │
│  • Missing episode numbers    • Logged, never blocks parsing                │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       EXTRACTOR (lexer/)                                    │
│  Walks AST, extracts quotes with metadata                                   │
│                                                                             │
│  ExtractedQuote {                                                           │
│      Content:     []DialogueElement                                         │
│      CharacterID: "27"                                                      │
│      AudioID:     "10100001"                                                │
│      Episode:     1                                                         │
│      Truth:       { HasRed: true, HasBlue: false, ... }                     │
│      SoundEffects: [{ SeNum: 47, AfterClip: 0 }, ...]                      │
│  }                                                                          │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                   TRANSFORMERS (lexer/transformer/)                         │
│  Converts raw AST to output formats                                         │
│                                                                             │
│  PlainTextTransformer ──► "Without love, it cannot be seen."                │
│  HtmlTransformer      ──► "<span class=\"red-truth\">...</span>"           │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ParsedQuote (dto/)                                  │
│                                                                             │
│  { Text, TextHtml, CharacterID, Character, AudioID, Episode,               │
│    SoundEffects, HasRedTruth, HasGoldTruth, ... }                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

`ParseFile` handles the decoding step before the pipeline above:

```
┌──────────────┐    ┌──────────────┐    ┌──────────────────────────────────┐
│  ONS2 .file  │───►│   Decoder    │───►│  Parse (pipeline above)          │
│  (encrypted) │    │  (decoder/)  │    │  returns quotes + subtitle refs  │
└──────────────┘    └──────────────┘    └──────────────────────────────────┘
```

Subtitle refs (pointers to `.ass` files found in the script) are returned alongside quotes for the caller to resolve.

## Installation

```bash
go get github.com/VictoriqueMoe/umineko_script_parser
```

## Usage

### From encrypted `.file` files

`ParseFile` takes an `io.Reader` for an ONS2-encrypted `.file`, decodes it (two-pass XOR + zlib), and parses the script into structured quotes.

```go
package main

import (
    "fmt"
    "os"

    scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func main() {
    f, err := os.Open("en.file")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    quotes, subtitleRefs, validationErrors, err := scriptparser.ParseFile(f)
    if err != nil {
        panic(err)
    }

    _ = subtitleRefs
    _ = validationErrors

    for _, q := range quotes {
        fmt.Printf("[EP%d] %s: %s\n", q.Episode, q.Character, q.Text)
    }
}
```

### From decoded script text

If you already have the raw script text (e.g. you decoded it yourself or are working with plain text exports), use `Parse` directly:

```go
quotes, subtitleRefs, validationErrors, err := scriptparser.ParseScriptText(rawScriptText)
if err != nil {
    log.Fatal(err)
}
```

## ParsedQuote

Each parsed quote contains:

| Field            | Type                | Description                                                                                |
|------------------|---------------------|--------------------------------------------------------------------------------------------|
| `Text`           | `string`            | Plain text content                                                                         |
| `TextHtml`       | `string`            | HTML with semantic markup (red/blue truth classes, ruby annotations, italic, colour spans) |
| `CharacterID`    | `string`            | Numeric character ID (e.g. `"10"` for Battler) or `"narrator"`                             |
| `Character`      | `string`            | Display name (e.g. `"Ushiromiya Battler"`)                                                 |
| `AudioID`        | `string`            | Comma-separated voice file IDs                                                             |
| `AudioCharMap`   | `map[string]string` | Audio ID to character ID mapping (multi-character quotes only)                             |
| `AudioTextMap`   | `map[string]string` | Audio ID to spoken text fragment (multi-audio quotes only)                                 |
| `Episode`        | `int`               | Episode number (1-8)                                                                       |
| `ContentType`    | `string`            | `""` (main story), `"tea"`, `"ura"`, or `"omake"`                                          |
| `HasRedTruth`    | `bool`              | Contains red truth                                                                         |
| `HasBlueTruth`   | `bool`              | Contains blue truth                                                                        |
| `HasGoldTruth`   | `bool`              | Contains gold truth                                                                        |
| `HasPurpleTruth` | `bool`              | Contains purple statements                                                                 |
| `SoundEffects`   | `[]SoundEffect`     | Associated sound effects with timing (`Filename`, `AfterClip`)                             |

`SoundEffect` has two fields: `Filename` (e.g. `"umise_047"`) and `AfterClip` (voice clip index the SE plays after, or `-1` for before all clips).

## HTML Output

The `TextHtml` field produces semantic HTML:

- Red truth: `<span class="red-truth">...</span>`
- Blue truth: `<span class="blue-truth">...</span>`
- Italic: `<em>...</em>`
- Colour: `<span style="color:#FF0000">...</span>`
- Ruby annotations: `<ruby>text<rp>(</rp><rt>reading</rt><rp>)</rp></ruby>`
- Line breaks: `<br>`

Dynamic preset colours from the script (gold text, purple text, etc.) are preserved as inline styles.

## Sub-packages

For advanced usage, the internals are fully exported:

| Package             | Description                                                          |
|---------------------|----------------------------------------------------------------------|
| `lexer`             | Tokenizer, recursive descent parser, quote extractor, validator      |
| `lexer/ast`         | AST node types (Script, DialogueLine, VoiceCommand, FormatTag, etc.) |
| `lexer/transformer` | Plain text and HTML transformers, preset context                     |
| `decoder`           | ONS2 format decryption (two-pass XOR + zlib)                         |
| `quote/character`   | 61 character constants with ID and name mappings                     |
| `quote/mutation`    | Post-parse correction engine (e.g. Kanon attribution fix)            |
| `dto`               | `ParsedQuote` type definition                                        |

### Working with the AST directly

```go
import (
    "github.com/VictoriqueMoe/umineko_script_parser/lexer"
    "github.com/VictoriqueMoe/umineko_script_parser/lexer/ast"
)

script := lexer.Parse(rawText)

for _, line := range script.Lines {
    switch l := line.(type) {
    case *ast.DialogueLine:
        voices := l.GetVoiceCommands()
    case *ast.EpisodeMarkerLine:
        fmt.Println("Episode:", l.Episode)
    }
}
```

### Custom transformers

```go
import "github.com/VictoriqueMoe/umineko_script_parser/lexer/transformer"

factory := transformer.NewFactory(transformer.NewPresetContext())
plain := factory.MustGet(transformer.FormatPlainText)
html := factory.MustGet(transformer.FormatHTML)

factory.Register(transformer.Format(100), myCustomTransformer)
```

## Supported Characters

All 61 voiced characters from Umineko are included with ID mappings:

Battler (`10`), Beatrice (`27`), Bernkastel (`28`), Lambdadelta (`29`), Erika (`46`), Dlanor (`47`), Featherine (`50`), and all family members, servants, witches, furniture, and Stakes of Purgatory.

```go
import "github.com/VictoriqueMoe/umineko_script_parser/quote/character"

name := character.CharacterNames.GetCharacterName(character.Battler)
// "Ushiromiya Battler"

ch := character.CharacterFromID("27")
// character.Beatrice
```

## License

MIT
