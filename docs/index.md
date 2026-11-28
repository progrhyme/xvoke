# xvoke Documentation

## Overview

xvoke standardizes cross-project work done by AI agents, and enforces quality gates on each change, so that the work is efficient and its quality is assured.
It is provided as a Claude Code plugin, and as a command-line tool, `xvoke`.

- **Project configuration**: each target project is described by a `project.yaml` file, which defines the working copy, the default branch, and the quality gates a change must pass.
- **Skill `xv-run`** (invoked as `/xvoke:xv-run`): carries out an instruction on a project, through the quality gates.
  By default, the skill works on a branch and takes the change to a merged pull request. A project can change this behavior in `project.yaml`.
- **Quality gates**: checks that run on the change before it is merged: `test`, `consistency`, `review` (by an AI agent or by the user), and `ci`.
- **Command `xvoke`**: validates a `project.yaml`, and writes the result as JSON. See [Installing the command](../README.md#installing-the-command).

The details are defined in the following files, which are the source of truth:

- [SKILL.md](../plugin/skills/xv-run/SKILL.md): the behavior of `xv-run`
- [project.example.yaml](../plugin/skills/xv-run/references/project.example.yaml): the format of `project.yaml`, with every field and gate kind
- [spec/cli.md](spec/cli.md): the interface of the `xvoke` command, and what `xvoke validate` checks

## Directory layout

`project.yaml` is placed in the workspace, the directory where the session runs. Each project has its own directory, `xvokedir/{project name}/`, and `/xvoke:xv-run {project name}` reads `project.yaml` from there. The working copy of a project can be anywhere; `root` in `project.yaml` points to it.

```
workspace/
└── xvokedir/
    └── myapp/            # the project name
        └── project.yaml

~/repos/myapp/            # the working copy
```

## Minimal project.yaml

```yaml
name: myapp
root: ~/repos/myapp
default_branch: main
quality_gates:
  - kind: test
    command: make test
  - kind: ci    # waits for the CI of the remote repository
```

## Usage

```text
/xvoke:xv-run myapp Add a --verbose option to the CLI
```

After you install the command (see [Installing the command](../README.md#installing-the-command)), you can check a `project.yaml` before a run:

```bash
xvoke validate xvokedir/myapp/project.yaml
```

## Contents

- [spec/](spec/): Specifications
- [adr/](adr/): Architecture Decision Records
- [develop/convention.md](develop/convention.md): Development conventions
- [maint/documents.md](maint/documents.md): Documentation conventions
