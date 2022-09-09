# gobreaking

I find breaking changes between two Go trees.

## Install

```shell
go install github.com/pierreprinetti/gobreaking
```

## Use

```shell
cd my/working/directory
gobreaking .
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
      - uses: pierreprinetti/gobreaking@v0.0.1
        id: gobreaking
        with:
          base: base
          new: new
      - name: Comment PR
        uses: thollander/actions-comment-pull-request@v1
        with:
          comment_includes: '## gobreaking hint'
          message: |
            ## gobreaking hint
            ${{ steps.gobreaking.outputs.verdict }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
