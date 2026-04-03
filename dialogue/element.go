package dialogue

type DialogueElement interface {
	ElementType() string
}

type TextElement interface {
	DialogueElement
	GetText() string
}

type ContainerElement interface {
	DialogueElement
	GetName() string
	GetParam() string
	GetContent() []DialogueElement
}

type SpecialCharElement interface {
	DialogueElement
	GetCharName() string
}
