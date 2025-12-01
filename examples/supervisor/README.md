# Supervisor Agent Example

This example demonstrates how to implement a multi-agent system using a Supervisor.

## Overview

In a Supervisor pattern, a central "supervisor" agent routes tasks to specialized worker agents. The workers perform their tasks and report back to the supervisor, which then decides the next step or finishes the workflow.

This example shows how to use `prebuilt.CreateSupervisor` to orchestrate multiple agents.

## Usage

```bash
go run main.go
```
