# Human-in-the-loop (HITL) Example

This example demonstrates how to implement a **Human-in-the-loop** workflow using LangGraphGo.

## 1. Background

In many real-world agentic applications, full automation is not always desirable or possible. You may need a human to:
- **Approve** critical actions (e.g., deploying code, sending emails).
- **Provide input** that the agent cannot access.
- **Correct** the agent's reasoning or state before it proceeds.

LangGraphGo supports this via an **Interrupt** mechanism. You can configure the graph to pause execution before or after specific nodes, allowing an external system (or human) to inspect and modify the state before resuming.

## 2. Key Concepts

- **InterruptBefore**: A configuration option that tells the graph to stop execution *before* entering a specified node.
- **GraphInterrupt**: A specific error type returned when the graph pauses. It contains the current state and the node where it stopped.
- **ResumeFrom**: A configuration option used to restart execution from a specific node, usually the one where it was interrupted.

## 3. How It Works

1.  **Define the Graph**: Standard graph definition with nodes and edges.
2.  **Initial Run**: Invoke the graph with `InterruptBefore` set to the target node (e.g., "human_approval").
3.  **Pause & Inspect**: The graph executes up to the target node and then returns a `GraphInterrupt` error. The current state is preserved.
4.  **Human Interaction**: The application catches the interrupt, presents the state to a human (simulated here), and updates the state based on their input (e.g., setting `Approved = true`).
5.  **Resume**: Invoke the graph again with the *updated state* and `ResumeFrom` set to the interrupted node. The graph continues execution from that point.

## 4. Code Highlights

### Setting up the Interrupt
```go
config := &graph.Config{
    InterruptBefore: []string{"human_approval"},
}
// The Invoke call will return a GraphInterrupt error when it hits "human_approval"
res, err := runnable.InvokeWithConfig(ctx, initialState, config)
```

### Handling the Interrupt
```go
var interrupt *graph.GraphInterrupt
if errors.As(err, &interrupt) {
    // Access the state at the moment of interruption
    currentState := interrupt.State.(State)
    // ... present to human ...
}
```

### Resuming Execution
```go
// Update state with human input
currentState.Approved = true 

// Resume configuration
resumeConfig := &graph.Config{
    ResumeFrom: []string{"human_approval"},
}
// Continue execution with the modified state
finalRes, err := runnable.InvokeWithConfig(ctx, currentState, resumeConfig)
```

## 5. Running the Example

```bash
go run main.go
```

**Expected Output:**
```text
=== Starting Workflow (Phase 1) ===
[Process] Processing request: Deploy to Production
Workflow interrupted at node: human_approval
Current State: {Input:Deploy to Production Approved:false Output:Processed: Deploy to Production}

=== Human Interaction ===
Reviewing request...
Approving request...

=== Resuming Workflow (Phase 2) ===
[Human] Request APPROVED.
[Finalize] Final output: Processed: Deploy to Production (Approved)
Workflow completed successfully.
Final Result: {Input:Deploy to Production Approved:true Output:Processed: Deploy to Production (Approved)}
```
