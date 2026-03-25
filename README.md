# umineko_script_parser

A Go library for parsing Umineko no Naku Koro ni (When the Seagulls Cry) game script files. Extracts structured dialogue quotes with character attribution, episode metadata, voice audio references, red/blue truth detection, and both plain text and HTML output.

## Installation

```bash
go get github.com/VictoriqueMoe/umineko_script_parser
```

## Usage

### From encrypted `.file` files

`NewLoader` handles the full pipeline: reads an ONS2-encrypted `.file` from the provided filesystem, decodes it (two-pass XOR + zlib), parses the script into structured quotes, resolves any ASS subtitle references, and applies post-parse corrections.

```go
package main

import (
    "embed"
    "fmt"

    scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

//go:embed data/*.file data/sub/*.ass
var dataFS embed.FS

func main() {
    loader := scriptparser.NewLoader(dataFS)
    quotes := loader.Load("en", "data/en.file")

    for _, q := range quotes {
        fmt.Printf("[EP%d] %s: %s\n", q.Episode, q.Character, q.Text)
    }
}
```

### From decoded script text

If you already have the raw script text (e.g. you decoded it yourself or are working with plain text exports), use `Parse` directly:

```go
quotes := scriptparser.Parse(rawScriptText)
```

## ParsedQuote

Each parsed quote contains:

| Field          | Type                | Description                                                                                |
|----------------|---------------------|--------------------------------------------------------------------------------------------|
| `Text`         | `string`            | Plain text content                                                                         |
| `TextHtml`     | `string`            | HTML with semantic markup (red/blue truth classes, ruby annotations, italic, colour spans) |
| `CharacterID`  | `string`            | Numeric character ID (e.g. `"10"` for Battler) or `"narrator"`                             |
| `Character`    | `string`            | Display name (e.g. `"Ushiromiya Battler"`)                                                 |
| `AudioID`      | `string`            | Comma-separated voice file IDs                                                             |
| `AudioCharMap` | `map[string]string` | Audio ID to character ID mapping (multi-character quotes only)                             |
| `AudioTextMap` | `map[string]string` | Audio ID to spoken text fragment (multi-audio quotes only)                                 |
| `Episode`      | `int`               | Episode number (1-8)                                                                       |
| `ContentType`  | `string`            | `""` (main story), `"tea"`, `"ura"`, or `"omake"`                                          |
| `HasRedTruth`  | `bool`              | Contains red truth                                                                         |
| `HasBlueTruth` | `bool`              | Contains blue truth                                                                        |

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
| `quote/loader`      | File loading with subtitle resolution and mutation pipeline          |
| `quote/mutation`    | Post-parse correction engine (e.g. Kanon attribution fix)            |
| `subtitle`          | ASS subtitle format parser                                           |
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
