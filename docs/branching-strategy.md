### Branching Strategy

- Work in this repository is done using feature-branching, with continous integration into the `dev` branch.
- A feature-branch should be created with reference to a specific GitHub issue.
- Before a pull-request is made, the developer of the feature-branch should verify that the feature works as intended and passes all tests (when implemented).
- Once work on a branch is completed, a pull-request into `dev` can be made, and the PR-template filled out.
- The PR-template should describe the what & why of the changes made, and reference the issue that it closes.
- A pull-request should be merged only after review from at least one other team member.

The general workflow is shown in the branch diagram below:

```mermaid
gitGraph
  commit id: "initial commit"
  branch dev
  checkout dev
  commit id: "dev setup"
  branch feature-1
  checkout feature-1
  commit id: "f1: commit 1"
  checkout dev
  branch feature-2
  checkout feature-2
  commit id: "f2: commit 1"
  checkout feature-1
  commit id: "f1: commit 2"
  checkout feature-2
  commit id: "f2: commit 2"
  checkout dev
  merge feature-1 id: "PR #1 (closes #10)"
  checkout feature-2
  commit id: "f2: commit 3"
  checkout dev
  merge feature-2 id: "PR #2 (closes #11)"
  checkout main
  merge dev tag: "v0.0.2"
```


