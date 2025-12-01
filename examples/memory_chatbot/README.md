# Memory Chatbot Example

This example demonstrates building a conversational chatbot with memory using LangGraphGo and LangChain memory integration.

## Overview

This chatbot maintains conversation context across multiple turns, allowing it to:
- Remember user information (name, preferences, etc.)
- Reference previous parts of the conversation
- Provide contextually relevant responses
- Support different memory strategies

## Features

- ✅ Full conversation memory with LangGraph integration
- ✅ Multiple memory types (Buffer, Window, Token)
- ✅ Automatic context management
- ✅ System message configuration
- ✅ Conversation history tracking
- ✅ Demo mode for testing without API keys

## Prerequisites

To run the full example with LLM integration, you need:

```bash
export OPENAI_API_KEY=your_api_key
# OR
export DEEPSEEK_API_KEY=your_api_key
```

If no API key is set, the example runs in demo mode with simulated responses.

## Usage

```bash
cd examples/memory_chatbot
go run main.go
```

## Architecture

The chatbot uses a graph-based architecture:

```
┌─────────────────────┐
│ add_system_message  │  Add system prompt (first turn only)
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│       chat          │  Load memory → LLM call → Save to memory
└─────────────────────┘
```

## Code Highlights

### Creating a Chatbot

```go
// Create chatbot with buffer memory
bot, err := NewChatbot(ctx, "buffer")

// Create chatbot with window memory (last 5 turns)
bot, err := NewChatbot(ctx, "window")

// Create chatbot with token buffer (1000 tokens max)
bot, err := NewChatbot(ctx, "token")
```

### Chatbot Implementation

```go
type Chatbot struct {
    llm    llms.Model
    memory *prebuilt.LangChainMemory
    graph  *graph.CompiledGraph[*ChatbotState]
}

func (c *Chatbot) Chat(ctx context.Context, message string) (string, error) {
    state := &ChatbotState{
        Input:    message,
        Messages: []llms.ChatMessage{},
    }
    
    result, err := c.graph.Invoke(ctx, state, nil)
    if err != nil {
        return "", err
    }
    
    return result.Output, nil
}
```

### Graph Node with Memory

```go
g.AddNode("chat", func(ctx context.Context, state *ChatbotState) (*ChatbotState, error) {
    // 1. Load conversation history from memory
    memVars, err := c.memory.LoadMemoryVariables(ctx, map[string]any{})
    if err != nil {
        return nil, err
    }
    
    // 2. Extract historical messages
    var historyMessages []llms.ChatMessage
    if history, ok := memVars["history"]; ok {
        if msgs, ok := history.([]llms.ChatMessage); ok {
            historyMessages = msgs
        }
    }
    
    // 3. Combine system + history + current input
    allMessages := append(state.Messages, historyMessages...)
    allMessages = append(allMessages, llms.HumanChatMessage{
        Content: state.Input,
    })
    
    // 4. Generate response
    response, err := c.llm.GenerateContent(ctx, allMessages)
    if err != nil {
        return nil, err
    }
    
    // 5. Save to memory
    err = c.memory.SaveContext(ctx, map[string]any{
        "input": state.Input,
    }, map[string]any{
        "output": response.Choices[0].Content,
    })
    
    state.Output = response.Choices[0].Content
    return state, nil
})
```

## Example Conversation

```
[Turn 1]
User: Hello! My name is Alice and I love programming in Go.
Bot: Hi Alice! It's great to meet you. Go is an excellent programming language!

[Turn 2]
User: What's my name?
Bot: Your name is Alice, as you mentioned in your first message.

[Turn 3]
User: What programming language do I like?
Bot: You mentioned that you love programming in Go!
```

## Memory Types

### ConversationBuffer
Stores all conversation history. Best for:
- Short conversations
- When full context is needed
- No token/memory constraints

### ConversationWindowBuffer
Keeps only the last N conversation turns. Best for:
- Long conversations
- When only recent context matters
- Memory-constrained environments

### ConversationTokenBuffer
Maintains conversation within a token limit. Best for:
- Token-aware applications
- Cost optimization
- API rate limit management

## API Reference

### Chatbot Methods

```go
// Create a new chatbot
func NewChatbot(ctx context.Context, memoryType string) (*Chatbot, error)

// Send a message and get response
func (c *Chatbot) Chat(ctx context.Context, message string) (string, error)

// Get conversation history
func (c *Chatbot) GetHistory(ctx context.Context) ([]llms.ChatMessage, error)

// Clear conversation history
func (c *Chatbot) ClearHistory(ctx context.Context) error
```

## Extending the Example

You can extend this chatbot by:

1. **Adding Tools**: Integrate function calling for actions
2. **Persistent Storage**: Save memory to database
3. **Multi-User Support**: Separate memory per user
4. **Streaming Responses**: Use streaming for real-time output
5. **Custom System Prompts**: Personalize bot behavior

## Learn More

- [Memory Basic Example](../memory_basic/)
- [LangGraphGo Documentation](../../docs/)
- [LangChain Memory Concepts](https://python.langchain.com/docs/modules/memory/)
