## gh-compare

This is a [`gh`](https://cli.github.com/) extension to replace the [`compare`](https://hub.github.com/hub-compare.1.html) command in the deprecated tool, [`hub`](https://hub.github.com/).

### Installation

```shell
gh extension install wassimk/gh-compare
```

### Usage

Compare the current branch against `main`. If this repository is a fork, it uses the upstream remote.

```shell
gh compare
```

Compare the current branch against another.

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
