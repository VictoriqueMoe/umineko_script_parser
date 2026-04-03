# When They Cry Script Parser

A Go library for parsing game script files from the When They Cry (07th Expansion) visual novel series. Supports both **Umineko no Naku Koro ni** and **Higurashi no Naku Koro ni**. Extracts structured dialogue quotes with character attribution, episode/arc metadata, voice audio references, sound effect associations, and both plain text and HTML output.

## Supported Games

| Game | Script Format | Arcs | Characters | Unique Features |
|------|--------------|------|------------|-----------------|
| Umineko | NScripter DSL with nested format tags | 8 episodes + tea/ura/omake | 61 | Red/blue/gold/purple truth detection, ruby annotations, semantic HTML |
| Higurashi | C-like function calls (OutputLine) | 19 arcs (8 main + 11 console) | 39+ named, 300+ total | Bilingual JP/EN text, arc identification |

## Installation

```bash
go get github.com/VictoriqueMoe/umineko_script_parser
```

## Quick Start

### Umineko

```go
import (
    "os"
    "github.com/VictoriqueMoe/umineko_script_parser/umineko"
)

f, _ := os.Open("en.file")
defer f.Close()

quotes, subtitleRefs, validationErrors, err := umineko.ParseFile(f)
for _, q := range quotes {
    fmt.Printf("[EP%d] %s: %s\n", q.Episode, q.Character, q.Text)
}
```

### Higurashi

```go
import (
    "os"
    "github.com/VictoriqueMoe/umineko_script_parser/higurashi"
)

f, _ := os.Open("en.file")
defer f.Close()

quotes, validationErrors, err := higurashi.ParseFile(f)
for _, q := range quotes {
    fmt.Printf("[%s] %s: %s\n", q.Arc, q.Character, q.Text)
    fmt.Printf("  JP: %s\n", q.TextJP)
}
```

Both games use the same ONS2 `.file` format (XOR + zlib encoded). Plain text input is also supported via `ParseScriptText()`.

## Output Types

### Shared Base (all games)

| Field | Type | Description |
|-------|------|-------------|
| `Text` | `string` | Plain text content |
| `TextHtml` | `string` | HTML-formatted text |
| `CharacterID` | `string` | Character identifier |
| `Character` | `string` | Display name |
| `AudioID` | `string` | Comma-separated voice file paths |
| `AudioCharMap` | `map[string]string` | Audio path to character ID mapping |
| `AudioTextMap` | `map[string]string` | Audio path to spoken text fragment |
| `Episode` | `int` | Episode/arc number |
| `ContentType` | `string` | Content type marker |
| `SoundEffects` | `[]SoundEffect` | Associated sound effects |

### Umineko-Specific (`dto.UminekoQuote`)

| Field | Type | Description |
|-------|------|-------------|
| `HasRedTruth` | `bool` | Contains red truth |
| `HasBlueTruth` | `bool` | Contains blue truth |
| `HasGoldTruth` | `bool` | Contains gold truth |
| `HasPurpleTruth` | `bool` | Contains purple statements |

### Higurashi-Specific (`dto.HigurashiQuote`)

| Field | Type | Description |
|-------|------|-------------|
| `TextJP` | `string` | Japanese plain text |
| `TextJPHtml` | `string` | Japanese HTML-formatted text |
| `Arc` | `string` | Arc name (e.g. `"onikakushi"`, `"watanagashi"`) |

## Architecture

### Shared Layer

```
dialogue/      DialogueElement interface + sub-interfaces (TextElement, ContainerElement, SpecialCharElement)
transformer/   Transformer interface, Factory, Format constants
dto/           BaseQuote (shared), UminekoQuote, HigurashiQuote
```

Both games implement `dialogue.DialogueElement` for their AST types, and `transformer.Transformer` for their output formatters. The shared `transformer.Factory` supports custom transformer registration for both games.

### Umineko Pipeline

```
ONS2 .file ──► Decoder ──► Lexer ──► Parser (AST) ──► Validator ──► Extractor ──► Transformers ──► UminekoQuote
```

The Umineko script uses a NScripter DSL with recursively nested format tags (`{p:1:{c:FF0000:{i:text}}}`), requiring a full tokenizer, recursive descent parser, and AST walker. The HTML transformer produces semantic markup with truth classes, ruby annotations, colour spans, and emphasis tags.

### Higurashi Pipeline

```
ONS2 .file ──► Decoder ──► Parser (line-by-line) ──► Transformers ──► HigurashiQuote
```

The Higurashi script uses flat C-like function calls (`OutputLine`, `ModPlayVoiceLS`, `ClearMessage`). A line-by-line state machine extracts quotes directly without tokenization or AST construction. Each `OutputLine` call contains both Japanese and English text.

## Sub-Packages

| Package | Description |
|---------|-------------|
| `dialogue` | Shared `DialogueElement` interface and sub-interfaces |
| `transformer` | Shared `Transformer` interface, `Factory`, `Format` constants |
| `dto` | `BaseQuote`, `UminekoQuote`, `HigurashiQuote`, `SoundEffect` |
| `umineko` | Umineko entry point (`ParseFile`, `ParseScriptText`) |
| `umineko/lexer` | Tokenizer, recursive descent parser, quote extractor, validator |
| `umineko/lexer/ast` | AST node types implementing `dialogue.DialogueElement` |
| `umineko/transformer` | Plain text and HTML transformers, preset context |
| `umineko/decoder` | ONS2 format decryption (two-pass XOR + zlib) |
| `umineko/character` | 61 character constants with ID and name mappings |
| `umineko/mutation` | Post-parse correction engine (e.g. Kanon attribution fix) |
| `higurashi` | Higurashi entry point (`ParseFile`, `ParseScriptText`) |
| `higurashi/ast` | Simple AST types implementing `dialogue.DialogueElement` |
| `higurashi/transformer` | Plain text and HTML transformers |
| `higurashi/character` | 39+ character constants with voice ID and name mappings |

### Custom Transformers

Both games support custom transformers via the shared factory:

```go
import (
    "github.com/VictoriqueMoe/umineko_script_parser/dialogue"
    "github.com/VictoriqueMoe/umineko_script_parser/transformer"
)

type MyTransformer struct{}

func (t *MyTransformer) Transform(elements []dialogue.DialogueElement) string {
    // custom output logic
}

factory := transformer.NewFactory()
factory.Register(transformer.Format(100), &MyTransformer{})
```

### Characters

```go
import umichar "github.com/VictoriqueMoe/umineko_script_parser/umineko/character"

name := umichar.CharacterNames.GetCharacterName(umichar.Battler)
// "Ushiromiya Battler"
```

```go
import higuchar "github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"

name := higuchar.CharacterNames.GetCharacterName(higuchar.Keiichi)
// "Maebara Keiichi"

ch := higuchar.CharacterFromName("Rena")
// higuchar.Rena
```

### Umineko AST

```go
import (
    "github.com/VictoriqueMoe/umineko_script_parser/umineko/lexer"
    "github.com/VictoriqueMoe/umineko_script_parser/umineko/lexer/ast"
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

## Higurashi Arc Order

The combined Higurashi script file contains arcs in this order:

**Main Arcs:** Onikakushi (1), Watanagashi (2), Tatarigoroshi (3), Himatsubushi (4), Meakashi (5), Tsumihoroboshi (6), Minagoroshi (7), Matsuribayashi (8)

**Console Arcs:** Someutsushi (9), Kageboshi (10), Tsukiotoshi (11), Taraimawashi (12), Yoigoshi (13), Tokihogushi (14), Miotsukushi Omote (15), Kakera (16), Miotsukushi Ura (17), Kotohogushi (18), Hajisarashi (19)

## License

MIT
