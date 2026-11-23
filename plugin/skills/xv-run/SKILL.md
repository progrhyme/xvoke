---
name: xv-run
description: Carries out an instruction on a registered project, from a work branch through quality gates to a merged pull request. Takes a project name and an instruction as its arguments. Runs only when invoked through the skill command (/xv-run); it has no natural-language trigger conditions.
---

# xv-run

Carry out the instruction on the project, and take the change through to a merged pull request.
The project's `project.yaml` defines where the working copy is and which quality gates the change must pass, so that the result meets the same standard on every project.

## Arguments

`/xv-run {project} {instruction}`: the first word is the project name, and the rest is the instruction. If either is missing, ask the user.

## Terms

- **AI review**: a gate of `kind: review` reviewed by an AI agent (`authorizer` omitted or `agent`).
- **User review**: a gate of `kind: review` with `authorizer: user`.

## Workflow

### 1. Load the project

- Read `references/project.example.yaml` first. It defines the format of `project.yaml` and what each field and gate kind means.
- Read `xvokedir/{project}/project.yaml` in the workspace (the directory the session started in). If it does not exist, tell the user and stop.
- Work in the directory `root`. Read the repository's own rules (`AGENTS.md`, `CLAUDE.md`, `CONTRIBUTING.md` and the like) before changing anything. They decide branch naming, commits, pushes, and pull requests.

### 2. Prepare a branch

- Branch from the latest `default_branch`, unless the user says otherwise. Work on the branch, never on `default_branch`, so that every change goes through a pull request.
- If the working copy has uncommitted changes you did not make, ask the user before going on.

### 3. Do the work

Do the work described by the instruction. Follow the repository's rules for style and structure.

### 4. Run the quality gates

Run every gate except `ci`, which needs the pull request and runs in step 6.

#### Order

- Run the gates in this order: `test`, `consistency`, AI review, user review. Use this order even if `project.yaml` lists them differently. Cheap mechanical checks go first, so that the user reviews a change that has already passed them.
- Gates of the same kind run in the listed order.

#### Selection

- The changed files are the difference from `default_branch`, including uncommitted and untracked files.
- Decide whether a gate runs by its `paths`, as `project.example.yaml` defines.

#### Execution

- Run each gate as its kind is defined in `project.example.yaml`. Fix what a gate reports, and run the gate again until it passes (for an AI review, follow the rule below). If you cannot make it pass, stop and report to the user.

#### AI review

- Give the subagent the gate's `description` and the files to review. The subagent does not know how the change was made, which keeps the review independent.
- After you fix the findings, ask the user whether to review again. Do not decide this yourself, because a repeated review costs time and the user knows how much the findings mattered.

#### User review

- Show what to review and the gate's `description`, and wait for the user's verdict.

### 5. Commit, push, and open the pull request

Follow the repository's rules on commits, pushes, and pull requests. If a rule is missing or ambiguous, ask the user instead of guessing, because these actions are hard to undo and visible to others.

### 6. Wait for the `ci` gate

Wait for the `ci` gate, if there is one, on the pull request.
If it fails, fix the cause. Before you push the fix, run again the gates of step 4 whose `paths` match the files the fix changed. Then push and wait again.

### 7. Merge the pull request

Merge the pull request by the repository's rules. If the rules do not say who merges or how, ask the user.

### 8. Report

Report the pull request, the gates that ran and their results, and anything left undone.
