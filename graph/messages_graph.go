package graph

// NewMessagesStateGraph creates a StateGraph with a default schema
// that handles "messages" using the AddMessages reducer.
// This is the recommended starting point for chat-based agents.
func NewMessagesStateGraph() *StateGraph {
	g := NewStateGraph()
	schema := NewMapSchema()
	schema.RegisterReducer("messages", AddMessages)
	g.SetSchema(schema)
	return g
}
