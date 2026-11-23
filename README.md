# xvoke

xvoke standardizes cross-project work done by AI agents, and enforces quality gates on each change, so that the work is efficient and its quality is assured.
It is provided as a Claude Code plugin.

Each target project is described by a `project.yaml` file at `xvokedir/{project-name}/project.yaml` in the workspace. The workspace is the directory where you run Claude Code, and it holds `xvokedir/` with one subdirectory per target project.
The `xv-run` skill reads `project.yaml` to take a change to the project through to a merged pull request.

## Installation

xvoke is installed as a plugin in a skills directory. At the start of each session, Claude Code loads every folder under `~/.claude/skills/` that has a `.claude-plugin/plugin.json` as a plugin.

1. Clone the repository, and move into it.

   ```bash
   git clone https://github.com/progrhyme/xvoke.git
   cd xvoke
   ```

2. Link the plugin directory into `~/.claude/skills/`.

   ```bash
   ln -s "$(pwd)/plugin" ~/.claude/skills/xvoke
   ```

3. Start a new session, or run `/reload-plugins` in a running session.

The skill is invoked as `/xvoke:xv-run`. To stop loading the plugin, run `claude plugin disable xvoke@skills-dir` (`skills-dir` marks a plugin loaded from a skills directory).

## Repository layout

```
AGENTS.md
CHANGELOG.md
LICENSE
README.md
docs/
├── index.md    # overview and usage
├── spec/       # Specifications
└── adr/        # Architecture Decision Records
plugin/         # the Claude Code plugin
├── .claude-plugin/plugin.json
└── skills/
    └── xv-run/
```

## Documentation

See [docs/index.md](docs/index.md).

## License

[MIT](LICENSE)
