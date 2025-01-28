# Devoid

Devoid (think de-void, like the Big Bang) creates an entire codebase from scratch, when provided with just a prompt and directory path.

It accomplishes this by using the configured LLM, by default the local `deepseek-r1:8b` model via Ollama. More models will be supported in the future.

Implementation is achieved via a state machine to interact with the LLM and process the results as structured output; think things like running bootstrap commands, creating directories and files, and writing code. The experience is guided in the terminal, with the ability to confirm LLM-driven actions, ask the LLM to modify the execution plan in arbitrary ways, answer clarifying questions for the LLM, and other shiny things.

## Current Status

The `initial` stage is implemented for project boostrapping, and the outputs (bootstrap commands and filesystem paths) are fed to the `scaffolding` stage. Next the scaffolding stage needs to actually use that output, and then call the next stage, which will write some code in the files created by the scaffolding stage.

Also on the roadmap is writing out checkpoint state (patches, model output and input, etc.) into the project repo for each stage to allow for easy rollback/restarting at a given stage.

## Demo

<video width="640" height="360" controls>
  <source src="./demo.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>
