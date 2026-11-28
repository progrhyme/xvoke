# xvoke

xvoke standardizes cross-project work done by AI agents, and enforces quality gates on each change, so that the work is efficient and its quality is assured.
It is provided as a Claude Code plugin, and as a command-line tool, `xvoke`, which validates a `project.yaml`.

Each target project is described by a `project.yaml` file at `xvokedir/{project-name}/project.yaml` in the workspace. The workspace is the directory where you run Claude Code, and it holds `xvokedir/` with one subdirectory per target project.
The `xv-run` skill reads `project.yaml` and carries out your instruction on the project through the quality gates. By default, the skill works on a branch and takes the change to a merged pull request; a project can opt out with `workflow_style.prefer_pull_request: false` in `project.yaml`.

## Installing the plugin

The xvoke plugin is installed in a skills directory. At the start of each session, Claude Code loads every folder under `~/.claude/skills/` that has a `.claude-plugin/plugin.json` as a plugin.

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

## The `xvoke` command

`xvoke validate` checks a `project.yaml` before a run, and finds a misspelled key or a missing required key. The command is optional: the `xv-run` skill does not use it at present.

### Installing the command

The archives of the command are attached to each [release](https://github.com/progrhyme/xvoke/releases), together with the checksums. In the names below, `{version}` is the version without the leading `v`.

| Platform | Archive |
| :--- | :--- |
| Linux, x86-64 | `xvoke_{version}_linux_amd64.tar.gz` |
| macOS, Intel | `xvoke_{version}_darwin_amd64.tar.gz` |
| macOS, Apple silicon | `xvoke_{version}_darwin_arm64.tar.gz` |
| Windows, x86-64 | `xvoke_{version}_windows_amd64.zip` |

An archive has a directory that holds the executable and `LICENSE`. The checksums of all archives of a release are in `xvoke_{version}_checksums.txt`.

This example installs version 0.0.3 on macOS with Apple silicon. Change the variables for another version or platform.

1. Download the archive and the checksums.

   ```bash
   version=0.0.3 os=darwin arch=arm64
   base="https://github.com/progrhyme/xvoke/releases/download/v${version}"
   curl -fLO "${base}/xvoke_${version}_${os}_${arch}.tar.gz"
   curl -fLO "${base}/xvoke_${version}_checksums.txt"
   ```

2. Verify the archive with `xvoke_{version}_checksums.txt`, for example with `sha256sum -c`, or with `shasum -a 256 -c` on macOS.

3. Extract the archive, and put the executable in a directory on `PATH`. This example uses `/usr/local/bin`, which needs `sudo`. Any directory on `PATH` works. If `/usr/local/bin` does not exist, create it first.

   ```bash
   tar -xzf "xvoke_${version}_${os}_${arch}.tar.gz"
   sudo install "xvoke_${version}_${os}_${arch}/xvoke" /usr/local/bin/
   ```

On Windows, extract the `.zip`, and put `xvoke.exe` from the extracted directory in a directory on `PATH`.

### Commands

| Command | What it does |
| :--- | :--- |
| `xvoke validate <project.yaml>` | Validates a `project.yaml`. |
| `xvoke schema` | Writes the latest JSON Schema of `project.yaml`. |
| `xvoke version` | Writes the version. |
| `xvoke help [<command>]` | Writes the usage. |

Every command writes one JSON document. The specification, including the output and the exit status, is in [docs/spec/cli.md](docs/spec/cli.md).

## Repository layout

```
AGENTS.md
CHANGELOG.md
LICENSE
Makefile        # build, test, and package the command
README.md
cmd/xvoke/      # entry point of the command
internal/       # Go packages of the command
docs/
├── index.md    # overview and usage
├── spec/       # Specifications
├── adr/        # Architecture Decision Records
└── develop/    # Development conventions
plugin/         # the Claude Code plugin
├── .claude-plugin/plugin.json
└── skills/
    └── xv-run/
.github/        # CI and release workflows
```

## Documentation

See [docs/index.md](docs/index.md).

## License

[MIT](LICENSE)
