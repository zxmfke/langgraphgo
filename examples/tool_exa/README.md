# Exa Search Agent Example

This example demonstrates how to use the Exa Search tool with a LangGraph ReAct Agent.

## Prerequisites

You need an Exa API key and an OpenAI-compatible API key (e.g., DeepSeek, OpenAI).

Set the environment variables:

```bash
export EXA_API_KEY=your_exa_key
export OPENAI_API_KEY=your_openai_key
# OR
export DEEPSEEK_API_KEY=your_deepseek_key
```

## Usage

Run the example:

```bash
go run main.go
```

## How it Works

1.  **Initialize LLM**: Connects to the LLM (defaulting to DeepSeek-V3).
2.  **Initialize Tool**: Creates the Exa search tool.
3.  **Create Agent**: Uses `prebuilt.CreateReactAgent` to build a graph-based agent that has access to the tool.
4.  **Execute**: The agent receives a query, decides to use the Exa tool to find information, and then synthesizes the answer.

## Code Overview

```go
// Create the agent with LLM and Tools
agent, err := prebuilt.CreateReactAgent(llm, []tools.Tool{exaTool})

// Invoke the agent
response, err := agent.Invoke(ctx, inputs)
```
