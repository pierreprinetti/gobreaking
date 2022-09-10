# gobreaking

I find breaking changes between two Go trees.

Given two versions of the same module or package, `gobreaking` finds breaking changes:
* Removing an exported name (constant, type, variable, function).
* Changing the type of an exported name.
* Adding or removing a method in an exported interface.
* Adding or removing a parameter in an exported function or interface.
* Changing the type of a parameter in an exported function or interface.
* Adding or removing a result in an exported function or interface.
* Changing the type of a result in an exported function or interface.
* Removing an exported field from an exported struct.
* Changing the type of an exported field of an exported struct.
* Adding an exported or unexported field to an exported struct containing only exported fields.
* Repositioning a field in an exported struct containing only exported fields.

Currently, only the first two rules are checked.

## Install

```shell
go install github.com/pierreprinetti/gobreaking@latest
```

## Use

```shell
gobreaking <path-to-base-module> <path-to-new-version>
```

## Use the Github action

```yaml
name: Find breaking changes
on:
  pull_request:
    types:
      - opened
      - synchronize
jobs:
  gobreaking:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          path: base
          ref: master
      - uses: actions/checkout@v3
        with:
          path: new
      - uses: pierreprinetti/gobreaking@v0.0.1   # or use a commit SHA
        id: gobreaking                           # id is needed to fetch the verdict in the comment step
        with:
          base: base
          new: new
      - name: Comment PR
        uses: thollander/actions-comment-pull-request@v1
        with:
          comment_includes: '## gobreaking hint' # Update the previous comment instead of creating a new one
          message: |
            ## gobreaking hint
            ${{ steps.gobreaking.outputs.verdict }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
