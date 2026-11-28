# 001. Add a Go command-line tool to validate project.yaml

| Item | Value |
| :--- | :--- |
| Status | ACCEPTED |
| Deciders | progrhyme |
| Updated | 2026-11-28 |

## Summary

- Context: the format of `project.yaml` is defined only by an annotated example that an AI agent reads, so a mistake in a `project.yaml` can go unnoticed until a run.
- Problem: how to find such mistakes mechanically before a run.
- Decision: add a command-line tool, `xvoke`, written in Go, that validates a `project.yaml`.

## Context and Problem Statement

xvoke is a Claude Code plugin. Each target project is described by a `project.yaml`, and the `xv-run` skill reads it to decide where the working copy is and which quality gates a change must pass.

The format of `project.yaml` is defined by `project.example.yaml`, an annotated example. The agent reads the example and interprets the file. Nothing that runs mechanically finds a mistake such as a misspelled key or a missing required key, so the mistake can go unnoticed until a run, or cause a wrong run.

The mistakes should be found before a run, by a validation that gives the same result each time. The validation must run on users' machines without a language runtime, and an agent must be able to read its output without parsing text.

## Decision Outcome

Chosen option: "a command-line tool written in Go", because it meets the drivers below, and a library for JSON Schema validation exists for Go.

The `xvoke` command validates a `project.yaml` against a JSON Schema, whose version can be selected, and writes the result as JSON, which an agent can read without parsing text. The validation is the first operation of the command. Other deterministic operations can be delegated to it in the future.

The decision drivers are:

- The validation gives the same result each time, and an agent can read the output without parsing text.
- It runs without a language runtime or dependencies on the user's machine, starts quickly, and uses little memory, so an agent can run it repeatedly.
- The code stays uniform and checkable: Go is statically typed, and has a small specification and standard tools.
- Other deterministic operations, whose result is determined by their input, can be delegated to it in the future.

## Pros and Cons of the Options

| Option | Pros | Cons |
| :--- | :--- | :--- |
| ✅️ A command-line tool written in Go | A single binary needs no runtime or dependencies for the user. A JSON Schema library exists. Statically typed, with a small specification and standard tools. | A program besides the plugin, which needs its own tests and lint. |
| Only the example and the prompt | Nothing to add to the repository. | Nothing checks the file mechanically, so the result depends on the agent's reading of the example. |
| A script in a language that needs a runtime, such as Python | Needs no build. | The user needs the runtime and the libraries. The versions on users' machines differ. The runtime starts on every run. |

## Consequences

### Good

- A misspelled key and a missing required key are found before a run, with the location of the mistake (a JSON Pointer and, when known, a line).
- The format of `project.yaml` has a schema, whose version is selected by the optional `schema_version` field.
- The command is a place to which more deterministic operations can be moved from the agent.

### Bad

- The user needs to install the command separately to use its functions.
- The repository has Go code besides the plugin, so it needs the tests and the lint of Go.

## More Information

- Other compiled languages were not compared.
- The specification of the `xvoke` command is in [spec/cli.md](../spec/cli.md).
- As of 2026-11-28, `xv-run` does not require `xvoke`.
