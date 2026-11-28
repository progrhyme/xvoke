# The xvoke command

`xvoke` is a command-line tool that is mainly run by an AI agent. This document specifies its interface: the commands, the output, the exit status, and what `xvoke validate` checks.

## Synopsis

```
xvoke <command> [flags]
```

| Command | What it does |
| :--- | :--- |
| `validate <project.yaml>` | Validates a `project.yaml`. |
| `schema` | Prints the latest JSON Schema of `project.yaml`. |
| `version` | Prints the version of the command. |
| `help [<command>]` | Prints the usage. |

## Output

- The command writes one JSON document to the standard output.
- The command writes nothing to the standard error output.
- The document is indented with two spaces, and ends with a newline.
- A command that could not run writes an error document, if the output can be written (see [Errors](#errors)).
- `schema` writes the schema itself, which is a JSON document.

## Exit status

| Status | Meaning |
| :--- | :--- |
| 0 | The command succeeded. For `validate`, the file has no errors. Warnings do not affect the status. |
| 1 | `validate` ran, and the file has errors. |
| 2 | The command could not run (see [Errors](#errors)). |

When writing the output fails, the status is 2, and the output may be incomplete.

## Commands

### validate

```
xvoke validate <project.yaml>
```

The argument is the path of one `project.yaml`. The command reads the file, and writes the result.

```json
{
  "valid": false,
  "path": "xvokedir/myapp/project.yaml",
  "errors": [
    {
      "source": "schema",
      "pointer": "/quality_gates/0",
      "line": 5,
      "message": "missing property 'command'"
    }
  ],
  "warnings": [
    {
      "source": "check",
      "pointer": "/root",
      "line": 2,
      "message": "directory does not exist: ~/repos/myapp"
    }
  ]
}
```

| Key | Meaning |
| :--- | :--- |
| `valid` | `true` if `errors` is empty. Warnings do not affect it. |
| `path` | The path as it was given in the argument. |
| `errors` | The problems that make the file invalid. |
| `warnings` | The problems that do not make the file invalid. |

`errors` and `warnings` are always arrays. They are empty arrays when there is nothing to report.

Each element of `errors` and `warnings` is a diagnostic:

| Key | Meaning |
| :--- | :--- |
| `source` | Where the problem was found: `schema` (a violation of the JSON Schema) or `check` (any other problem that the code finds, such as invalid YAML, a bad `schema_version`, or a failed check of the code). |
| `pointer` | A [JSON Pointer](https://www.rfc-editor.org/rfc/rfc6901) to the place of the problem. It is the empty string for the whole document. |
| `line` | The line in the YAML. The key is omitted when the line is unknown. |
| `message` | An English description, for people and agents to read. Its wording is not part of the interface. |

#### What is checked

The command reads the YAML, selects the schema by `schema_version`, validates against it, and then makes the checks that the schema cannot express.

- **YAML.** The file holds at most one YAML document. A syntax error, or a second document, is an error with the `source` `check`, and the other steps are skipped.
- **Schema version.** `schema_version` is an optional integer that selects the schema. When it is omitted, the latest version is used. When it is not an integer, or no schema has that version, the error has the `source` `check` and the `pointer` `/schema_version`, and the schema is not applied. How a new version is made is in [convention.md](../develop/convention.md).
- **Schema.** `xvoke schema` writes the schema, which is the definition of the format of `project.yaml`. A violation of the schema is an error with the `source` `schema`. The schema rejects a key that it does not define, so that a misspelled key is reported, at the key itself. An unknown key in a gate is reported only when the `kind` of the gate is valid. An empty file, or a document that is not an object, is also a violation, and has no `line`.
- **Checks of the code.** The command makes these checks on `name` and `root`, each only when it is a non-empty string, even if the schema reports errors or is not applied:

  | Check | Level | `pointer` |
  | :--- | :--- | :--- |
  | `name` equals the name of the directory that holds the file. For `xvokedir/myapp/project.yaml`, it is `myapp`. | Error | `/name` |
  | `root` is an existing directory. `~`, or a path that starts with `~/`, is relative to the home directory. Another relative path is relative to the current directory of the command. | Warning | `/root` |

In `errors`, the error of the schema version, or the diagnostics of the schema, come first. The diagnostics of the schema are sorted by `line`, then `pointer`, then `message`, and a diagnostic that has no line comes first. The errors of the checks of the code follow. `warnings` has only the warnings of the checks of the code.

### schema

```
xvoke schema
```

The command writes the latest JSON Schema of `project.yaml`, in JSON. It takes no argument. The keys of an object in the schema are currently sorted alphabetically, but the order is not part of the interface.

### version

```
xvoke version
```

```json
{
  "version": "0.0.3"
}
```

The command takes no argument. The version is the module version that Go records when the command is built.

- A binary built from a tagged commit has the version of the tag, without the leading `v`.
- A build from another commit records a pseudo-version, such as `0.0.3-0.20261128042310-6e8106f3ae53`.
- `+dirty` is added to the version when the working tree has changes.
- The version is `dev` when the build records no module version, as with `go run` and `go test`.

### help

```
xvoke help [<command>]
```

`xvoke help` writes the usage and the list of the commands.

```json
{
  "usage": "xvoke <command> [flags]",
  "commands": [
    {
      "name": "validate",
      "summary": "Validate a project.yaml. ...",
      "usage": "xvoke validate <project.yaml>",
      "flags": [
        {
          "names": ["-h", "--help"],
          "description": "Show the usage of the command."
        }
      ]
    }
  ]
}
```

`xvoke help <command>` writes one command, in the same form as an element of `commands`, under the key `command`:

```json
{
  "command": {
    "name": "validate",
    "summary": "...",
    "usage": "xvoke validate <project.yaml>",
    "flags": [{ "names": ["-h", "--help"], "description": "Show the usage of the command." }]
  }
}
```

Every command lists the help flag, after the flags that the command defines.

These forms write the help, and the status is 0:

| Form | Writes |
| :--- | :--- |
| `xvoke --help`, `xvoke -h`, `xvoke help` | The commands |
| `xvoke help <command>` | The command |
| `xvoke <command> --help`, `xvoke <command> -h` | The command |

After the name of a known command, `-h` or `--help` anywhere before `--` prints the help of that command, regardless of the other arguments. A flag after `--` is not a flag: `xvoke validate -- -h` reads a file named `-h`. Other spellings, such as `-help`, are not part of the interface.

## Errors

A command that could not run writes an error document, and the status is 2.

```json
{
  "error": {
    "message": "cannot read the file: open project.yaml: no such file or directory"
  }
}
```

`message` is an English description. The cases are:

- No command is given.
- The command is unknown, or a flag is given in place of the command, as in `xvoke --foo`.
- The arguments do not fit the command: a flag that the command does not define, or a wrong number of arguments. `validate` takes one argument, `schema` and `version` take none, and `help` takes at most one.
- The file cannot be read.
- The validation cannot run, for example because the embedded schema cannot be compiled.
