## gh-compare

This is a [`gh`](https://cli.github.com/) extension to replace the [`compare`](https://hub.github.com/hub-compare.1.html) command in the mostly deprecated tool, [`hub`](https://hub.github.com/).

### Installation

```shell
gh extension install wassimk/gh-compare
```

### Usage

Compare the current branch against the `main` branch.

```shell
gh compare
```

Compare the current branch against another branch.

```shell
gh compare other_branch
```

Compare any two branches.


```shell
gh compare their_branch...my_branch
```

Compare any two commits.

```shell
gh compare 7a67154..205b073
```
