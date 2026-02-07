#!/bin/bash
timestamp=$(date -u +%Y%m%d%H%M%S)
task_name=${1:-task}
filename="${timestamp}_${task_name}.md"
cat > "agent-work/$filename" << 'TEMPLATE'
# <Task Name>

## Status: in_progress

## Context
<!-- What problem is being solved, why is it needed -->

## Value Proposition
<!-- What the feature achieves, acceptance criteria -->

## Alternatives considered
<!-- Other approaches and why this was chosen -->

## Todos
- [ ] Todo 1
- [ ] Todo 2

## Notes
TEMPLATE
echo "Created work file: agent-work/$filename"
