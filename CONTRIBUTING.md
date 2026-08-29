# Contributing to toolbox-go

We would love for you to contribute to `toolbox-go` and help make it even better than it is today!
As a contributor, we would like you to follow the guidelines.

## Coding Rules

To ensure consistency throughout the source code, keep these rules in mind as you are working:

* All features or bug fixes **must be tested** by one or more unit-tests.
* All exported API items **must be documented** with godoc comments, starting with the name of the
  item being documented.
* `golangci-lint run` and `go test ./...` **must pass** before opening a pull request. The linter
  configuration covers vetting as well as formatting, and wraps all code at **100 characters**.
* `go fix ./...` **must leave the code unchanged**: whatever it rewrites is applied before opening
  the pull request, not left to the reviewer.
* Whatever the project builds or runs **must be reproducible**: pin what can be pinned, from the
  module dependencies in `go.sum` to the actions of a workflow, so that the same revision still
  gives the same result tomorrow. See [reproducible builds][reproducible-builds] for the rationale.
* The files configuring the repository (`.gitignore`, editor and tooling configuration, ...)
  describe the project, **not a workstation**. Anything that only makes sense on a single machine
  or for a single setup belongs to the per user equivalent the tool provides:
  `~/.config/git/ignore` for git ignore rules, user level settings for editors, and so on.
* The repository is a **single Go module**. Each library is a package grouped by domain (`auth`,
  `crypto`, `session`, ...) and keeps its public API as small as possible: anything that does not
  need to be exported stays unexported.
* Code shared between libraries that is not part of the public API belongs to `internal/`.
* An adapter under `adapters/` stays thin and **must not leak its third-party package** outside of
  its own directory: it only translates a library API into what the third party expects. Anything
  more than translation belongs to the library itself.
* Prefer the standard library. Adding a third-party dependency must be justified in the pull
  request: as everything lives in a single module, every consumer inherits it.
* All packages share the same module version: any **breaking change** to an exported API must be
  called out in the pull request, as it forces a major version bump for the whole module.

## Commit Message Guidelines

We have very precise rules over how our git commit messages can be formatted. This leads to **more
readable messages** that are easy to follow when looking through the **project history**.

### Commit Message Format

As we follow the [Conventional Commits][conventional-commits] specification, each commit message
consists of a **header**, a **body** and a **footer**. The header has a special format that includes
a **type**, a **scope** and a **subject**:

```text
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory and the **scope** of the header is optional.

Any line of the commit message cannot be longer 100 characters! This allows the message to be easier
to read on GitHub as well as in various git tools.

The footer should contain a [closing reference to an issue][github-issue-closing] if any.

Samples:

```text
docs(readme): correct spelling
```

```text
fix(session): stop handing back expired sessions

The store only dropped an entry when its cleanup ran, so a session read between two sweeps was
returned long after its expiry. Check the expiry on read as well.

Refs: #123
```

### Revert

If the commit reverts a previous commit, it should begin with `revert:`, followed by the header of
the reverted commit. In the body it should say: `This reverts commit <hash>.`, where the hash is the
SHA of the commit being reverted.

### Type

Should be one of the following:

* **build**: Changes that affect the build system or external dependencies
* **chore**: Maintenance tasks that don't affect the delivered code or the build (e.g. tooling
  configuration, release commits)
* **ci**: Changes to our CI configuration files and scripts
* **docs**: Documentation only changes
* **feat**: A new feature
* **fix**: A bug fix
* **perf**: A code change that improves performance
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **revert**: Reverts a previous commit (see the [Revert](#revert) section)
* **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing
  semi-colons, etc)
* **test**: Adding missing tests or correcting existing tests

### Scope

The scope should be the root package the change is about, as perceived by the person reading the
project history: the name of the top level directory the code lives in.

Going one level deeper, as `<root package>/<subpackage>`, is allowed when that subpackage name is,
and will remain, unique in the repository: the reader then knows exactly what was touched, and the
scope stays valid as the layout grows. Fall back to the root package alone as soon as the change
spans several subpackages, or the subpackage name could later be reused somewhere else.

Not every change belongs to a package. Rather than an exhaustive list of exceptions, use common
sense and pick whatever helps the person reading the history:

* a `docs` change about the repository itself rather than about a package (the README, this file,
  any other top level document) takes no scope
* a `build` change concerning a single target may use that target as its scope, be it an operating
  system or an architecture (e.g. `build(windows)`, `build(arm64)`)
* a change spread over every package takes no scope, as there is nothing useful to name
  (e.g. `style: apply golangci-lint fmt`)

The scope is optional: leaving it out is always better than inventing one.

### Subject

The subject contains a succinct description of the change:

* use the imperative, present tense: "change" not "changed" nor "changes"
* don't capitalize the first letter
* no dot (.) at the end

### Body

Just as in the **subject**, use the imperative, present tense: "change" not "changed" nor "changes".
The body should include the motivation for the change and contrast this with previous behavior.

### Footer

The footer should contain any information about **Breaking Changes** and is also the place to
reference GitHub issues that this commit **Closes**.

**Breaking Changes** should start with the word `BREAKING CHANGE:` with a space or two newlines. The
rest of the commit message is then used for this.

## License

This project is licensed under either of the Apache License, Version 2.0
([LICENSE-APACHE](LICENSE-APACHE)) or the MIT license ([LICENSE-MIT](LICENSE-MIT)), at your option.

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in
the work by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any
additional terms or conditions.

[conventional-commits]: https://www.conventionalcommits.org/en/v1.0.0/
[github-issue-closing]: https://help.github.com/articles/closing-issues-via-commit-messages/
[reproducible-builds]: https://reproducible-builds.org/
