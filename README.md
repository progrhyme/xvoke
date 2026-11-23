# xvoke

xvoke standardizes cross-project work so that it can be carried out efficiently with guaranteed quality.
It is designed for work done by AI agents, and is provided as a Claude Code plugin.

Each target project is described by a `project.yaml` file in the workspace (`xvokedir/{project-name}/project.yaml`).
The `xv-run` skill uses it to work on the project and take a change to the project through to a merged pull request.

## Repository layout

```
AGENTS.md
README.md
docs/
├── index.md    # overview and usage
└── adr/        # Architecture Decision Records
plugin/         # the Claude Code plugin
├── .claude-plugin/plugin.json
└── skills/
    └── xv-run/
```

## Documentation

See [docs/index.md](docs/index.md).
