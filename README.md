# Devoid

Devoid (think de-void, like the Big Bang) creates an entire codebase from scratch, when provided with just a prompt and directory path.

It accomplishes this by using the configured LLM, by default the local `deepseek-r1:8b` model via Ollama. Any Ollama model is supported, and more models will be supported in the future.

Implementation is achieved via a state machine to interact with the LLM and process the results as structured output; think things like running bootstrap commands, creating directories and files, and writing code. The experience is guided in the terminal, with the ability to confirm LLM-driven actions, ask the LLM to modify the execution plan in arbitrary ways, answer clarifying questions for the LLM, and other shiny things.

## Safety

All LLM outputs are processed through safety and other validations before moving to the next stage. For things that can't reasonably be validated like arbitrary commands to run, a warning is displayed next to the list of actions so the user can personally validate them before proceeding. Interactive safety checks can be dangerously skipped with `--skip-interactive-safety-checks`.

## Current Status

The `initial` stage is implemented for project boostrapping, and the outputs are fed to the `ast` stage. Next the AST stage needs to use that output to create a directed graph / adjacency list, and then call the next stage, which will write those files to disk. Code will follow after.

Also on the roadmap is writing out checkpoint state (patches, model output and input, etc.) into the project repo for each stage to allow for easy rollback/restarting at a given stage.

## Demo

![](./demo.gif)
