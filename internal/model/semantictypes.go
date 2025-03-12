package model

type SemanticType int

const (
	ClassType SemanticType = iota
	PropertyType
	StringType
	CommentType
)

func (s SemanticType) String() string {
	switch s {
	case ClassType:
		return "class"
	case PropertyType:
		return "property"
	case StringType:
		return "string"
	case CommentType:
		return "comment"
	default:
		return "unknown"
	}
}

func SemanticTypes() []SemanticType {
	return []SemanticType{ClassType, PropertyType, StringType, CommentType}
}

func SemantincTypeAsString(s SemanticType) string {
	return s.String()
}

type TokenReady struct {
	Line   int
	Column int
	Length int
	Type   SemanticType
}
