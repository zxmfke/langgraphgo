# Tools Package

This package provides implementations of various tools that adhere to the `langchaingo` Tool interface.

## Available Tools

### Tavily Search

Uses the [Tavily API](https://tavily.com/) for web search.

**Usage:**

```go
import "github.com/smallnest/langgraphgo/tool"

// Create a new Tavily search tool
// It will look for TAVILY_API_KEY environment variable if apiKey is empty
tavilyTool, err := tool.NewTavilySearch("", tool.WithTavilySearchDepth("advanced"))
if err != nil {
    log.Fatal(err)
}

// Use the tool
result, err := tavilyTool.Call(context.Background(), "what is langgraphgo?")
```

### Exa Search

Uses the [Exa API](https://exa.ai/) for neural search.

**Usage:**

```go
import "github.com/smallnest/langgraphgo/tool"

// Create a new Exa search tool
// It will look for EXA_API_KEY environment variable if apiKey is empty
exaTool, err := tool.NewExaSearch("", tool.WithExaNumResults(10))
if err != nil {
    log.Fatal(err)
}

// Use the tool
result, err := exaTool.Call(context.Background(), "golang agent frameworks")
```

## Interface

All tools implement the following interface:

```go
type Tool interface {
    Name() string
    Description() string
    Call(ctx context.Context, input string) (string, error)
}
```
