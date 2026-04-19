# When They Cry Script Parser

A Go library for parsing game script files from the When They Cry (07th Expansion) visual novel series. Supports **Umineko no Naku Koro ni**, **Higurashi no Naku Koro ni**, and **Ciconia no Naku Koro ni** (Phase 1). Extracts structured dialogue quotes with character attribution, episode/arc/chapter metadata, voice audio references (where present), sound effect associations, and both plain text and HTML output.

## Supported Games

| Game      | Script Format                                   | Taxonomy                                                   | Characters                               | Unique Features                                                                                  |
|-----------|-------------------------------------------------|------------------------------------------------------------|------------------------------------------|--------------------------------------------------------------------------------------------------|
| Umineko   | NScripter DSL with nested format tags           | 8 episodes + tea/ura/omake                                 | 61                                       | Red/blue/gold/purple truth detection, ruby annotations, semantic HTML                            |
| Higurashi | C-like function calls (OutputLine)              | 19 arcs (8 main + 11 console)                              | 39+ named, 300+ total                    | Bilingual JP/EN text, arc identification                                                         |
| Ciconia   | Ponscripter-fork-wh dialect (`langen`/`langjp`) | Prologue + 25 acts + act25b + Epilogue + 16 data fragments | 148 script markers, 49 curated main cast | Inline hex color spans, Keropoyo color-based attribution, synthetic content-hash IDs (voiceless) |

## Installation

```bash
go get github.com/VictoriqueMoe/umineko_script_parser
```

## Quick Start

All three games share the same shape: instantiate a `Parser` for the game, feed it to the generic `scriptparser.ParseReader` (decoded file) or `scriptparser.ParseText` (plain string).

### Umineko

```go
import (
    "os"
    scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
    "github.com/VictoriqueMoe/umineko_script_parser/umineko"
)

f, _ := os.Open("en.file")
defer f.Close()

parser := umineko.NewParser()
quotes, validationErrors, err := scriptparser.ParseReader(f, parser)
subtitleRefs := parser.SubtitleRefs()  // umineko-specific side output

for _, q := range quotes {
    fmt.Printf("[EP%d] %s: %s\n", q.Episode, q.Character, q.Text)
}
```

### Higurashi

```go
import (
    "os"
    scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
    "github.com/VictoriqueMoe/umineko_script_parser/higurashi"
)

f, _ := os.Open("en.file")
defer f.Close()

quotes, validationErrors, err := scriptparser.ParseReader(f, higurashi.NewParser())
for _, q := range quotes {
    fmt.Printf("[%s] %s: %s\n", q.Arc, q.Character, q.Text)
    fmt.Printf("  JP: %s\n", q.TextJP)
}
```

### Ciconia

```go
import (
    "os"
    scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
    "github.com/VictoriqueMoe/umineko_script_parser/ciconia"
)

f, _ := os.Open("ciconia_en.file")
defer f.Close()

quotes, validationErrors, err := scriptparser.ParseReader(f, ciconia.NewParser())
for _, q := range quotes {
    fmt.Printf("[%s/%s] %s: %s\n", q.ContentType, q.Chapter, q.Character, q.Text)
    fmt.Printf("  JP: %s\n", q.TextJP)
}
```

All three games use the same ONS2 `.file` format (XOR + zlib encoded) when read via `ParseReader`. If you already have plain script text, use `scriptparser.ParseText(script, parser)` instead - it skips the decode step. To produce a `.file` from plain text, use the included `cmd/encoder` tool: `go run ./cmd/encoder input.txt output.file`.

## Output Types

### Shared Base (all games)

| Field          | Type                | Description                        |
|----------------|---------------------|------------------------------------|
| `Text`         | `string`            | Plain text content                 |
| `TextHtml`     | `string`            | HTML-formatted text                |
| `CharacterID`  | `string`            | Character identifier               |
| `Character`    | `string`            | Display name                       |
| `AudioID`      | `string`            | Comma-separated voice file paths   |
| `AudioCharMap` | `map[string]string` | Audio path to character ID mapping |
| `AudioTextMap` | `map[string]string` | Audio path to spoken text fragment |
| `Episode`      | `int`               | Episode/arc number                 |
| `ContentType`  | `string`            | Content type marker                |
| `SoundEffects` | `[]SoundEffect`     | Associated sound effects           |

### Umineko-Specific (`dto.UminekoQuote`)

| Field            | Type   | Description                |
|------------------|--------|----------------------------|
| `HasRedTruth`    | `bool` | Contains red truth         |
| `HasBlueTruth`   | `bool` | Contains blue truth        |
| `HasGoldTruth`   | `bool` | Contains gold truth        |
| `HasPurpleTruth` | `bool` | Contains purple statements |

### Higurashi-Specific (`dto.HigurashiQuote`)

| Field        | Type     | Description                                     |
|--------------|----------|-------------------------------------------------|
| `TextJP`     | `string` | Japanese plain text                             |
| `TextJPHtml` | `string` | Japanese HTML-formatted text                    |
| `Arc`        | `string` | Arc name (e.g. `"onikakushi"`, `"watanagashi"`) |

### Ciconia-Specific (`dto.CicroniaQuote`)

| Field        | Type     | Description                                                                                                  |
|--------------|----------|--------------------------------------------------------------------------------------------------------------|
| `TextJP`     | `string` | Japanese plain text                                                                                          |
| `TextJPHtml` | `string` | Japanese HTML-formatted text                                                                                 |
| `Chapter`    | `string` | Chapter id: `"00"` (prologue), `"01"`–`"25"`, `"25b"`, `"ep"` (epilogue), `"df01"`–`"df16"` (data fragments) |

`ContentType` is one of `"prologue"`, `"chapter"`, `"epilogue"`, or `"data_fragment"`. `Episode` is `1` (Phase 1) for all Ciconia quotes. Because Ciconia is voiceless, `AudioID` is a synthetic stable hash (`pro:xxxxxxxx`, `c14:xxxxxxxx`, `ep:xxxxxxxx`, `df03:xxxxxxxx`, with `:N` suffix on within-chapter hash collisions) derived from character + EN + JP text - it survives line reordering within a chapter. Inline hex color spans are rendered in `TextHtml` as `<span style="color:#RRGGBB">...</span>`. Lines that open in Keropoyo's distinctive green `#8df270` are automatically attributed to `character.Keropoyo` even though the script lacks an explicit speaker marker for them.

## Architecture

### Shared Layer

```
dialogue/      DialogueElement interface + sub-interfaces (TextElement, ContainerElement, SpecialCharElement)
transformer/   Transformer interface, Factory, Format constants
dto/           BaseQuote (shared), UminekoQuote, HigurashiQuote, CicroniaQuote
```

All three games implement `dialogue.DialogueElement` for their AST types, and `transformer.Transformer` for their output formatters. The shared `transformer.Factory` supports custom transformer registration for every game.

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

### Ciconia Pipeline

```
Plain script ──► Parser (line-by-line) ──► Body tokenizer ──► Synthetic ID ──► Transformers ──► CicroniaQuote
```

Ciconia uses the Ponscripter-fork-wh dialect (same engine as Umineko, unrelated dialogue syntax). Pairs of `langjp【Name】` / `langen^Name:^` mark speakers; `langjp「...」\` and `langen^"..."^\` carry bilingual dialogue. A line-by-line state machine detects chapter labels (`*prologue`, `*act1`–`*act25`, `*act25b`, `*tips001`–`*tips016`) and collects consecutive `lang*` lines into quotes. The epilogue is split out from `*act25b` at the `movie "movie\p1last\..."` cutscene trigger (the engine's "phase 1 last" marker) into a synthetic `"ep"` chapter. A secondary tokenizer strips `^@^` pauses, `^!wNNN^` timing hints, and quote-wrap `^` markers while preserving `^#RRGGBB` color spans as `ColoredText` AST nodes. Because the script has no audio IDs, a stable `AudioID` is synthesized per-quote by hashing `(characterId + EN + JP)` under a chapter prefix.

## Sub-Packages

| Package                 | Description                                                                                                                  |
|-------------------------|------------------------------------------------------------------------------------------------------------------------------|
| `dialogue`              | Shared `DialogueElement` interface and sub-interfaces                                                                        |
| `transformer`           | Shared `Transformer` interface, `Factory`, `Format` constants                                                                |
| `dto`                   | `BaseQuote`, `UminekoQuote`, `HigurashiQuote`, `CicroniaQuote`, `SoundEffect`                                                |
| `.` (root)              | Generic `ParseText[Q]` / `ParseReader[Q]` helpers, `Parser[Q]` interface, `ValidationError`, `ErrBinaryInput`                |
| `umineko`               | Umineko parser (`umineko.NewParser()`), plus `SubtitleRefs()` accessor                                                       |
| `umineko/lexer`         | Tokenizer, recursive descent parser, quote extractor, validator                                                              |
| `umineko/lexer/ast`     | AST node types implementing `dialogue.DialogueElement`                                                                       |
| `umineko/transformer`   | Plain text and HTML transformers, preset context                                                                             |
| `decoder`               | ONS2 format decryption + encoding (two-pass XOR + zlib), shared by all games                                                 |
| `umineko/character`     | 61 character constants with ID and name mappings                                                                             |
| `umineko/mutation`      | Post-parse correction engine (e.g. Kanon attribution fix)                                                                    |
| `higurashi`             | Higurashi parser (`higurashi.NewParser()`)                                                                                   |
| `higurashi/ast`         | Simple AST types implementing `dialogue.DialogueElement`                                                                     |
| `higurashi/transformer` | Plain text and HTML transformers                                                                                             |
| `higurashi/character`   | 39+ character constants with voice ID and name mappings                                                                      |
| `ciconia`               | Ciconia parser (`ciconia.NewParser()`)                                                                                       |
| `ciconia/ast`           | `PlainText` and `ColoredText` AST types                                                                                      |
| `ciconia/transformer`   | Plain text and HTML transformers (with inline color span rendering)                                                          |
| `ciconia/character`     | 148 script-marker enum constants (`CharacterNames`) + a curated 49-entry `MainCharacterNames` subset for main-cast filtering |
| `cmd/encoder`           | CLI tool to encode a plain script text file into the ONS2 `.file` format                                                     |

### Custom Transformers

All three games support custom transformers via the shared factory:

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

```go
import cichar "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"

// Display name - Ciconia names are the single-word form exactly as they
// appear in the script markers (no fabricated full names):
name := cichar.CharacterNames.GetCharacterName(cichar.Miyao)
// "Miyao"

ch := cichar.CharacterFromName("Jayden")
// cichar.Jayden

// Main-cast subset (49 named speakers, excluding ensemble/role labels
// like "AOU Officer", "Reporter", "Announcer"). Use this to build
// "curated + additional" filter UIs:
main := cichar.MainCharacterNames.GetAllCharacters()
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

## Ciconia Chapter Layout

The Ciconia Phase 1 script is organized into four content types keyed by the `Chapter` field:

**Prologue:** `"00"` (~69 quotes)

**Main Acts:** `"01"`–`"25"` + `"25b"` (Ciconia's final act is split into `*act25` and `*act25b` script labels)

**Epilogue:** `"ep"` - automatically split from the tail of `*act25b` at the `movie "movie\p1last\..."` cutscene trigger (the engine's "phase 1 last" marker). ~39 quotes, starting with Koshka's "......I did it..." soliloquy.

**Data Fragments:** `"df01"`–`"df16"` (the 16 menu-accessible tips entries)

Typical parse of the Phase 1 full script yields ~9,350 quotes across 44 chapter units, with ~390 quotes containing inline hex color spans and ~347 attributed to Keropoyo via the green-color heuristic.

## License

MIT
