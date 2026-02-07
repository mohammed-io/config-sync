#!/bin/bash
task_name=${1:-task}
for f in agent-work/*_${task_name}.md; do
  if [ -f "$f" ]; then
    basename=$(basename "$f")
    sed -i '' 's/Status: in_progress/Status: completed '"$(date -u +%Y%m%d%H%M%S)"'/' "$f" 2>/dev/null || \
    sed -i 's/Status: in_progress/Status: completed '"$(date -u +%Y%m%d%H%M%S)"'/' "$f"
    mv "$f" "agent-work/completed/$basename"
    echo "Completed and moved: agent-work/completed/$basename"
  fi
done
