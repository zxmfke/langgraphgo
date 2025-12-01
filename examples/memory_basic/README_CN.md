# Memory 基础示例

本示例演示如何将 LangChain 的 memory 功能集成到 LangGraphGo 中。

## 概述

Memory 集成允许您在多次交互中维护对话上下文。本示例展示了各种 memory 类型和使用模式：

1. **ConversationBuffer**: 存储所有对话历史
2. **ConversationWindowBuffer**: 仅保留最后 N 轮对话
3. **ConversationTokenBuffer**: 在 token 限制内维护历史
4. **ChatMessageHistory**: 直接的消息历史管理
5. **自定义 Memory 键**: 使用自定义的输入/输出/memory 键

## 功能特性

- ✅ 多种 memory 类型（Buffer、Window、Token）
- ✅ 与 LangGraph 集成
- ✅ 自定义 memory 配置
- ✅ 消息历史管理
- ✅ 跨轮次保持上下文

## 使用方法

```bash
cd examples/memory_basic
go run main.go
```

## 代码要点

### 基础 ConversationBuffer

```go
// 创建对话缓冲 memory
mem := prebuilt.NewConversationBufferMemory(
    memory.WithReturnMessages(true),
)

// 保存对话上下文
mem.SaveContext(ctx, map[string]any{
    "input": "你好，我叫 Alice",
}, map[string]any{
    "output": "你好 Alice！很高兴认识你。",
})

// 获取所有消息
messages, _ := mem.GetMessages(ctx)
```

### ConversationWindowBuffer

```go
// 仅保留最后 2 轮对话
mem := prebuilt.NewConversationWindowBufferMemory(2,
    memory.WithReturnMessages(true),
)
```

### ChatMessageHistory

```go
history := prebuilt.NewChatMessageHistory()

// 添加不同类型的消息
history.AddMessage(ctx, llms.SystemChatMessage{
    Content: "你是一个有帮助的助手。",
})
history.AddUserMessage(ctx, "你好！")
history.AddAIMessage(ctx, "你好！有什么可以帮助你的吗？")

// 获取所有消息
messages, _ := history.Messages(ctx)
```

### 自定义 Memory 键

```go
mem := prebuilt.NewConversationBufferMemory(
    memory.WithInputKey("user_input"),
    memory.WithOutputKey("ai_response"),
    memory.WithMemoryKey("chat_history"),
    memory.WithReturnMessages(true),
)
```

### 与 LangGraph 集成

```go
// 创建 memory
mem := prebuilt.NewConversationBufferMemory(
    memory.WithReturnMessages(true),
)

// 在图节点中使用
g.AddNode("chat", func(ctx context.Context, state *State) (*State, error) {
    // 加载 memory
    memVars, _ := mem.LoadMemoryVariables(ctx, map[string]any{})
    
    // 获取历史消息
    var historyMessages []llms.ChatMessage
    if history, ok := memVars["history"]; ok {
        if msgs, ok := history.([]llms.ChatMessage); ok {
            historyMessages = msgs
        }
    }
    
    // 在 LLM 调用中使用历史
    allMessages := append(historyMessages, llms.HumanChatMessage{
        Content: state.Input,
    })
    
    response, _ := llm.GenerateContent(ctx, allMessages)
    
    // 保存到 memory
    mem.SaveContext(ctx, map[string]any{
        "input": state.Input,
    }, map[string]any{
        "output": response.Choices[0].Content,
    })
    
    return state, nil
})
```

## Memory 类型对比

| Memory 类型              | 描述            | 使用场景                 |
| ------------------------ | --------------- | ------------------------ |
| ConversationBuffer       | 存储所有消息    | 短对话，需要完整上下文   |
| ConversationWindowBuffer | 保留最后 N 轮   | 长对话，仅需最近上下文   |
| ConversationTokenBuffer  | 维护 token 限制 | Token 感知应用，成本控制 |

## API 参考

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
    // 方法
    AddMessage(ctx context.Context, message llms.ChatMessage) error
    AddUserMessage(ctx context.Context, message string) error
    AddAIMessage(ctx context.Context, message string) error
    Messages(ctx context.Context) ([]llms.ChatMessage, error)
    Clear(ctx context.Context) error
    SetMessages(ctx context.Context, messages []llms.ChatMessage) error
}
```

## 了解更多

- [LangChain Memory 文档](https://python.langchain.com/docs/modules/memory/)
- [LangGraphGo 文档](../../docs/)
- [Memory 聊天机器人示例](../memory_chatbot/)
