# Documentation conventions

Conventions for writing the documents in `docs/`. For the file layout, see [index.md](../index.md).

## General

- The primary sources of truth for the specification are the schema and the source code. Documents are secondary information.

## Specifications (`docs/spec/`)

A specification has the following properties. Keep them in mind when writing.

- It helps readers understand a complex specification.
  - For a specification that carries a lot of information and is hard to express in the schema or in source code comments, manage the specification document as the source of truth.
- It is a reference for some features.
  - For example, CLI commands. Describe the API accurately, following the format.
- It does not cover the whole specification.
  - Copying the primary sources item by item only raises the maintenance cost.
