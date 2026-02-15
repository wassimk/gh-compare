## gh-compare

A [`gh`](https://cli.github.com/) CLI extension that opens GitHub's compare view in your browser. It detects forks, resolves your repository's default branch automatically, and supports branch and commit comparisons.

This was built to replace the [`compare`](https://hub.github.com/hub-compare.1.html) command from the deprecated [`hub`](https://hub.github.com/) tool.

### 🛠️ Installation

```shell
gh extension install wassimk/gh-compare
```

### ✨ Usage

The command must be run inside a git repository with a GitHub remote.

Compare the current branch against the repository's default branch.

```shell
gh compare
```

Compare the current branch against another branch.

```shell
gh compare other_branch
```

Compare any two branches using `...` (three-dot) notation. This shows only the changes introduced by the second branch, like a pull request diff. Commits added to the first branch after the two diverged are ignored.

```shell
gh compare main...my_branch
```

You can also use `..` (two-dot) notation for a direct diff between two branches as they exist right now. Unlike three-dot, this includes the effect of commits on both sides.

```shell
gh compare main..my_branch
```

Compare any two commits.

```shell
gh compare 7a67154..205b073
```

### 🍴 Forks

If your repository is a fork with an `upstream` remote, `gh compare` automatically generates a cross-fork compare URL. For example, running `gh compare` on a fork will produce a URL like:

```
https://github.com/upstream-owner/repo/compare/main...your-user:your-branch
```

This matches GitHub's compare-across-forks behavior so the diff is shown against the upstream repository.

### 📦 Upgrade / Uninstall

```shell
gh extension upgrade wassimk/gh-compare
gh extension remove gh-compare
```
