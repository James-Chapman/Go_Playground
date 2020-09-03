package graph

type Edge struct {
	vertex Vertex // Originating vertex
	to Vertex
	cost int64
}

