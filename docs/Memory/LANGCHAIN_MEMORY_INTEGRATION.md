# LangChain Memory Integration

LangGraphGo provides seamless integration with [LangChainGo's memory package](https://github.com/tmc/langchaingo/tree/main/memory), allowing you to easily manage conversation history and context in your graph-based applications.

## Overview

Memory is crucial for conversational AI applications. It allows the system to:
- Remember previous interactions
- Maintain context across multiple turns
- Personalize responses based on history

We provide a `LangChainMemory` adapter that wraps LangChainGo's memory implementations, making them compatible with LangGraphGo's architecture.

## Available Memory Types

The integration supports all standard LangChainGo memory types:

### 1. ConversationBuffer
Stores the entire conversation history. This is useful when you want the model to have access to everything that has been said.

```go
mem := prebuilt.NewConversationBufferMemory(
    memory.WithReturnMessages(true),
)
```

### 2. ConversationWindowBuffer
Keeps a sliding window of the most recent N conversation turns. This is useful for keeping the context size manageable while retaining recent context.

```go
// Keep only the last 5 turns
mem := prebuilt.NewConversationWindowBufferMemory(5,
    memory.WithReturnMessages(true),
)
```

### 3. ConversationTokenBuffer
Maintains conversation history within a specific token limit. This is ideal for managing costs and staying within LLM context window limits.

```go
// Keep history within 1000 tokens
mem := prebuilt.NewConversationTokenBufferMemory(llm, 1000,
    memory.WithReturnMessages(true),
)
```

### 4. ChatMessageHistory
Provides direct access to the underlying message history storage, allowing you to manually add or retrieve messages.

```go
history := prebuilt.NewChatMessageHistory()
history.AddUserMessage(ctx, "Hello")
history.AddAIMessage(ctx, "Hi there")
```

## Integration with LangGraph

Integrating memory into a LangGraph node typically involves three steps:

1.  **Load**: Retrieve conversation history from memory.
2.  **Process**: Use the history (combined with new input) to generate a response.
3.  **Save**: Store the new input and output back into memory.

### Example Pattern

```go
// 1. Define your graph node
func chatNode(ctx context.Context, state *State) (*State, error) {
    // --- LOAD ---
    // Load memory variables (e.g., history)
    memVars, _ := mem.LoadMemoryVariables(ctx, map[string]any{})
    
    // Extract messages from the "history" key
    var historyMessages []llms.ChatMessage
    if history, ok := memVars["history"]; ok {
        if msgs, ok := history.([]llms.ChatMessage); ok {
            historyMessages = msgs
        }
    }
    
    // --- PROCESS ---
    // Combine history with the current user input
    allMessages := append(historyMessages, llms.HumanChatMessage{
        Content: state.Input,
    })
    
    // Call the LLM
    response, _ := llm.GenerateContent(ctx, allMessages)
    aiResponse := response.Choices[0].Content
    
    // --- SAVE ---
    // Save the turn to memory
    mem.SaveContext(ctx, map[string]any{
        "input": state.Input,
    }, map[string]any{
        "output": aiResponse,
    })
    
    // Update state
    state.Output = aiResponse
    return state, nil
}
```

## Customization

### Custom Keys
You can customize the keys used for input, output, and memory variables if your application requires specific naming conventions.

```go
mem := prebuilt.NewConversationBufferMemory(
    memory.WithInputKey("user_query"),      // Default: "input"
    memory.WithOutputKey("bot_response"),   // Default: "output"
    memory.WithMemoryKey("chat_log"),       // Default: "history"
    memory.WithReturnMessages(true),
)
```

### Human/AI Prefixes
You can also customize the prefixes used for messages when they are returned as a string (i.e., `WithReturnMessages(false)`).

```go
mem := prebuilt.NewConversationBufferMemory(
    memory.WithHumanPrefix("User"),
    memory.WithAIPrefix("Assistant"),
)
```

## Examples

We provide complete examples demonstrating these concepts:

- **[Basic Memory Examples](../../examples/memory_basic/)**: Shows how to use different memory types and configurations.
- **[Chatbot with Memory](../../examples/memory_chatbot/)**: A complete chatbot simulation demonstrating state management and memory persistence.
