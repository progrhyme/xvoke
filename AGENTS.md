# AGENTS.md

Rules for AI agents that edit files, commit, or open pull requests in this repository.

## Overview

xvoke is a Claude Code plugin that standardizes cross-project work. It ships a command-line tool, also named `xvoke`, written in Go. See [README.md](README.md).

## Language

- Write all files, commit messages, and pull request text in English.

## Commits

- Start the subject with an emoji and a prefix, then a short description: `✨ feat: add xv-run skill`.
  - ✨ `feat:` new feature
  - 🐛 `fix:` bug fix
  - 📝 `docs:` documentation
  - 🎨 `style:` formatting only
  - ♻️ `refactor:`
  - ⚡️ `perf:`
  - ✅ `test:`
  - 🏗️ `build:` build system or dependencies
  - 👷 `ci:` CI configuration or scripts
  - 🔧 `chore:` anything else
  - ⏪️ `revert:`
- Make each commit one kind of change. Split changes with different prefixes into separate commits, unless one cannot stand without the other (e.g. a feature and its tests).
- Keep every commit in a consistent, working state.
- If a commit contains several related changes (e.g. a feature and its tests), list them in the body.

## Pull requests

- Do not push to `main` directly. Work on a branch and merge it through a pull request.
- Write the PR title in the same format as a commit subject.
- Describe what changed and why in the PR description.

## Layout

| Path | Contents |
| :--- | :--- |
| `docs/spec/` | Specifications |
| `docs/adr/` | Architecture Decision Records |
| `docs/develop/` | Development conventions |
| `plugin/skills/` | Skills provided by the plugin (e.g. `xv-run`) |
| `cmd/xvoke/` | Entry point of the `xvoke` command |
| `internal/` | Go packages of the `xvoke` command |
| `Makefile` | Targets to build, check, and package the command |
| `.github/` | CI and release workflows and their scripts |
