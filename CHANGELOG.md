# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

Each entry starts with a label that tells what it changes:

- `(skill)`: the behavior and the format of the skills
- `(plugin)`: the plugin package
- `(go)`: the `xvoke` command, written in Go, and the JSON Schema it embeds
- `(doc)`: the documentation
- `(repo)`: the repository itself

## [Unreleased]

### Added

- (skill) The optional `schema_version` key of `project.yaml`.
- (go) The Go module of the `xvoke` command, and the JSON Schema of `project.yaml` that it embeds.
- (go) The `xvoke` command with `validate`, `schema`, `version`, and `help`. It writes JSON to the standard output.

## [0.0.2] - 2026-11-23

### Added

- (skill) The `workflow_style.prefer_pull_request` setting of `project.yaml`. When it is `false`, `xv-run` works on the default branch directly and opens no pull request.

## [0.0.1] - 2026-11-23

### Added

- (skill) The `xv-run` skill, which carries out an instruction on a project through to a merged pull request.
- (skill) The `project.yaml` configuration file of a project, with its quality gates.
- (plugin) The Claude Code plugin manifest.
- (doc) The overview, the usage, and the installation guide.
- (repo) The MIT license.

[Unreleased]: https://github.com/progrhyme/xvoke/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/progrhyme/xvoke/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/progrhyme/xvoke/releases/tag/v0.0.1
