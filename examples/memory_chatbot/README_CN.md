# Memory 聊天机器人示例

本示例演示如何使用 LangGraphGo 和 LangChain memory 集成构建具有记忆功能的对话机器人。

## 概述

这个聊天机器人在多轮对话中维护对话上下文，使其能够：
- 记住用户信息（姓名、偏好等）
- 引用对话的先前部分
- 提供上下文相关的响应
- 支持不同的记忆策略

## 功能特性

- ✅ 完整的对话记忆与 LangGraph 集成
- ✅ 多种 memory 类型（Buffer、Window、Token）
- ✅ 自动上下文管理
- ✅ 系统消息配置
- ✅ 对话历史跟踪
- ✅ 无需 API 密钥的演示模式

## 前置条件

要运行完整的 LLM 集成示例，您需要：

```bash
export OPENAI_API_KEY=your_api_key
# 或者
export DEEPSEEK_API_KEY=your_api_key
```

如果未设置 API 密钥，示例将以演示模式运行，使用模拟响应。

## 使用方法

```bash
cd examples/memory_chatbot
go run main.go
```

## 架构

聊天机器人使用基于图的架构：

```
┌─────────────────────┐
│ add_system_message  │  添加系统提示（仅首轮）
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│       chat          │  加载记忆 → LLM 调用 → 保存到记忆
└─────────────────────┘
```

## 代码要点

### 创建聊天机器人

```go
// 使用缓冲 memory 创建聊天机器人
bot, err := NewChatbot(ctx, "buffer")

// 使用窗口 memory 创建（保留最后 5 轮）
bot, err := NewChatbot(ctx, "window")

// 使用 token 缓冲创建（最多 1000 tokens）
bot, err := NewChatbot(ctx, "token")
```

### 聊天机器人实现

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

### 带 Memory 的图节点

```go
g.AddNode("chat", func(ctx context.Context, state *ChatbotState) (*ChatbotState, error) {
    // 1. 从 memory 加载对话历史
    memVars, err := c.memory.LoadMemoryVariables(ctx, map[string]any{})
    if err != nil {
        return nil, err
    }
    
    // 2. 提取历史消息
    var historyMessages []llms.ChatMessage
    if history, ok := memVars["history"]; ok {
        if msgs, ok := history.([]llms.ChatMessage); ok {
            historyMessages = msgs
        }
    }
    
    // 3. 组合系统消息 + 历史 + 当前输入
    allMessages := append(state.Messages, historyMessages...)
    allMessages = append(allMessages, llms.HumanChatMessage{
        Content: state.Input,
    })
    
    // 4. 生成响应
    response, err := c.llm.GenerateContent(ctx, allMessages)
    if err != nil {
        return nil, err
    }
    
    // 5. 保存到 memory
    err = c.memory.SaveContext(ctx, map[string]any{
        "input": state.Input,
    }, map[string]any{
        "output": response.Choices[0].Content,
    })
    
    state.Output = response.Choices[0].Content
    return state, nil
})
```

## 示例对话

```
[第 1 轮]
用户: 你好！我叫 Alice，我喜欢用 Go 编程。
机器人: 你好 Alice！很高兴认识你。Go 是一门优秀的编程语言！

[第 2 轮]
用户: 我叫什么名字？
机器人: 你的名字是 Alice，正如你在第一条消息中提到的。

[第 3 轮]
用户: 我喜欢什么编程语言？
机器人: 你提到你喜欢用 Go 编程！
```

## Memory 类型

### ConversationBuffer
存储所有对话历史。最适合：
- 短对话
- 需要完整上下文时
- 无 token/内存限制

### ConversationWindowBuffer
仅保留最后 N 轮对话。最适合：
- 长对话
- 仅需最近上下文时
- 内存受限环境

### ConversationTokenBuffer
在 token 限制内维护对话。最适合：
- Token 感知应用
- 成本优化
- API 速率限制管理

## API 参考

### Chatbot 方法

```go
// 创建新的聊天机器人
func NewChatbot(ctx context.Context, memoryType string) (*Chatbot, error)

// 发送消息并获取响应
func (c *Chatbot) Chat(ctx context.Context, message string) (string, error)

// 获取对话历史
func (c *Chatbot) GetHistory(ctx context.Context) ([]llms.ChatMessage, error)

// 清除对话历史
func (c *Chatbot) ClearHistory(ctx context.Context) error
```

## 扩展示例

您可以通过以下方式扩展此聊天机器人：

1. **添加工具**: 集成函数调用以执行操作
2. **持久化存储**: 将 memory 保存到数据库
3. **多用户支持**: 为每个用户分离 memory
4. **流式响应**: 使用流式传输实现实时输出
5. **自定义系统提示**: 个性化机器人行为

## 了解更多

- [Memory 基础示例](../memory_basic/)
- [LangGraphGo 文档](../../docs/)
- [LangChain Memory 概念](https://python.langchain.com/docs/modules/memory/)
