# Memory Basic Example

This example demonstrates how to integrate LangChain's memory functionality with LangGraphGo.

## Overview

The memory integration allows you to maintain conversation context across multiple interactions. This example shows various memory types and usage patterns:

1. **ConversationBuffer**: Stores all conversation history
2. **ConversationWindowBuffer**: Keeps only the last N conversation turns
3. **ConversationTokenBuffer**: Maintains history within a token limit
4. **ChatMessageHistory**: Direct message history management
5. **Custom Memory Keys**: Using custom input/output/memory keys

## Features

- ✅ Multiple memory types (Buffer, Window, Token)
- ✅ Integration with LangGraph
- ✅ Custom memory configuration
- ✅ Message history management
- ✅ Context preservation across turns

## Usage

```bash
cd examples/memory_basic
go run main.go
```

## Code Highlights

### Basic ConversationBuffer

```go
// Create a conversation buffer memory
mem := prebuilt.NewConversationBufferMemory(
    memory.WithReturnMessages(true),
)

// Save conversation context
mem.SaveContext(ctx, map[string]any{
    "input": "Hello, my name is Alice",
}, map[string]any{
    "output": "Hi Alice! Nice to meet you.",
})

// Get all messages
messages, _ := mem.GetMessages(ctx)
```

### ConversationWindowBuffer

```go
// Keep only the last 2 conversation turns
mem := prebuilt.NewConversationWindowBufferMemory(2,
    memory.WithReturnMessages(true),
)
```

### ChatMessageHistory

```go
history := prebuilt.NewChatMessageHistory()

// Add different types of messages
history.AddMessage(ctx, llms.SystemChatMessage{
    Content: "You are a helpful assistant.",
})
history.AddUserMessage(ctx, "Hello!")
history.AddAIMessage(ctx, "Hi! How can I help you?")

// Get all messages
messages, _ := history.Messages(ctx)
```

### Custom Memory Keys

```go
mem := prebuilt.NewConversationBufferMemory(
    memory.WithInputKey("user_input"),
    memory.WithOutputKey("ai_response"),
    memory.WithMemoryKey("chat_history"),
    memory.WithReturnMessages(true),
)
```

### Integration with LangGraph

```go
// Create memory
mem := prebuilt.NewConversationBufferMemory(
    memory.WithReturnMessages(true),
)

// In your graph node
g.AddNode("chat", func(ctx context.Context, state *State) (*State, error) {
    // Load memory
    memVars, _ := mem.LoadMemoryVariables(ctx, map[string]any{})
    
    // Get historical messages
    var historyMessages []llms.ChatMessage
    if history, ok := memVars["history"]; ok {
        if msgs, ok := history.([]llms.ChatMessage); ok {
            historyMessages = msgs
        }
    }
    
    // Use history in your LLM call
    allMessages := append(historyMessages, llms.HumanChatMessage{
        Content: state.Input,
    })
    
    response, _ := llm.GenerateContent(ctx, allMessages)
    
    // Save to memory
    mem.SaveContext(ctx, map[string]any{
        "input": state.Input,
    }, map[string]any{
        "output": response.Choices[0].Content,
    })
    
    return state, nil
})
```

## Memory Types Comparison

| Memory Type              | Description           | Use Case                                 |
| ------------------------ | --------------------- | ---------------------------------------- |
| ConversationBuffer       | Stores all messages   | Short conversations, full context needed |
| ConversationWindowBuffer | Keeps last N turns    | Long conversations, recent context only  |
| ConversationTokenBuffer  | Maintains token limit | Token-aware applications, cost control   |

## API Reference

### LangChainMemory

```go
type Memory interface {
    SaveContext(ctx context.Context, inputValues, outputValues map[string]any) error
    LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error)
    Clear(ctx context.Context) error
    GetMessages(ctx context.Context) ([]llms.ChatMessage, error)
}
```

### ChatMessageHistory

```go
type ChatMessageHistory struct {
    // Methods
    AddMessage(ctx context.Context, message llms.ChatMessage) error
    AddUserMessage(ctx context.Context, message string) error
    AddAIMessage(ctx context.Context, message string) error
    Messages(ctx context.Context) ([]llms.ChatMessage, error)
    Clear(ctx context.Context) error
    SetMessages(ctx context.Context, messages []llms.ChatMessage) error
}
```

## Learn More

- [LangChain Memory Documentation](https://python.langchain.com/docs/modules/memory/)
- [LangGraphGo Documentation](../../docs/)
- [Memory Chatbot Example](../memory_chatbot/)
