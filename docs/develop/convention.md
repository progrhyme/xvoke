# Development conventions

This document sets the conventions for developing xvoke. It currently covers how the versions of the schema of `project.yaml` are made.

## Versions of the schema

The schema of `project.yaml` is versioned. A `project.yaml` can name the version that it follows with `schema_version`; see [the specification](../spec/cli.md) for how the command uses it.

- The maintainers make a new version only for a change that makes a file that was valid invalid. Adding an optional key does not make a new version: they add the key to the latest version.
- Every version keeps its schema file in `internal/schema/`, named `project.schema.v{N}.yaml`. A file can therefore keep following an older version.
- [project.example.yaml](../../plugin/skills/xv-run/references/project.example.yaml) describes only the latest version.
