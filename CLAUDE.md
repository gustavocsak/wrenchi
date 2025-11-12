# Claude-Code instructions for wrenchi TUI

## Context
- Language: Go
- Project: `wrenchi` – Linux systems-diagnostic tool
- Current state: CLI already prints RAM usage; CPU parsing proc/stat and proc/cpuinfo done
- Goal: add a real-time, interactive TUI that shows both RAM and CPU stats

## Claude Goals
- Help the user plan and design the tool, the user wants to learn
- Do not implement code unless told otherwise, give tips and assistance to the user
- Only implement code when the user requests it, usually when it is something tedious like replacing function names
