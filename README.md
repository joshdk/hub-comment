[![CircleCI][circleci-badge]][circleci-link]
[![License][license-badge]][license-link]
[![Github downloads][github-downloads-badge]][github-release-link]
[![GitHub release][github-release-badge]][github-release-link]

# Hub-Comment

📝 Automate GitHub pull request comments in CI

## Installing

### From source

A development version can be installed by running:

```bash
$ go get -u github.com/joshdk/hub-comment
```

### Precompiled binary

A prebuilt [release][github-release-link] binary can be installed by running:

```bash
$ wget -q https://github.com/joshdk/hub-comment/releases/download/0.1.0/hub-comment_linux_amd64
$ sudo install hub-comment_linux_amd64 /usr/bin/hub-comment
```

## Usage

### Environment Setup

Since `hub-comment` communicates with the [GitHub API](https://developer.github.com/v3/), you must first create and export an API token into the working environment.

```bash
$ export GITHUB_TOKEN='2b6c...f4bd'
```

### Comment Template

You can write a comment template file, using the same syntax used for [Go templates](https://golang.org/pkg/text/template/). For example, you could save the following as `hello-template.txt`:

```
🤖 Hello from CI!

Build #{{.Build.Number}} was run successfully.
```

### Running

```
$ hub-comment -template-file hello-template.txt
Posting new comment as Hub Bot (hub-bot) on joshdk/hub-comment#123:

→ 🤖 Hello from CI!
→
→ Build #123 was run successfully.

To view comment visit:
→ https://github.com/joshdk/hub-comment/issues/123#issuecomment-421483151
```

### Templating

At its core, `hub-comment` is a templating engine that applies context about the current build and PR, to a text template, producing a comment body. Templates are interpreted by Go's [`text/template` package](https://golang.org/pkg/text/template/). Several pieces of relevant context are made available to the template.

#### Build Context

The `.Build` object contains values that are derived from the current CircleCI build context. More information on this environment can be found [in the CircleCI docs](https://circleci.com/docs/2.0/env-vars/#built-in-environment-variables).

For example, the following template could be rendered as so:

```
[Build #{{.Build.Number}}]({{.Build.URL}}) is now complete!
```

```
[Build #123](https://circleci.com/gh/joshdk/hub-comment/123) is now complete!
```

Template | Content | Rendered example
---|---|---
`{{.Build.CI}}`| Are we running in CI right now? | `true`
`{{.Build.Index}}`| Index of the specific build instance. |`0`
`{{.Build.Job}}`| Name of the current job. |`deploy`
`{{.Build.Nodes}}`| Number of total build instances. |`1`
`{{.Build.Number}}`| Number of the CircleCI build. |`123`
`{{.Build.Owner}}`| User or org that owns the current project. |`joshdk`
`{{.Build.Repo}}`| Name of the repository of the current project. |`hub-comment`
`{{.Build.Stage}}`| Name of the current job. |`deploy`
`{{.Build.URL}}`| URL for the current build. |`https://circleci.com/gh/joshdk/hub-comment/123`
`{{.Build.User}}`| GitHub username of the user who triggered the build. |`joshdk`
`{{.Build.Workflow}}`| Unique identifier for the workflow instance of the current job. |`c7417d05-f6a0-4a6b-af5d-f3248b584d3f`

#### Env Context

The `.Env` object contains the current working environment. Useful for passing in custom data for use in templating (version strings, urls, error messages, etc).

For example, the following template could be rendered as so:

```
Deployed version v{{.Env.snapshot}} to staging!
```

```
Deployed version v1.2.3-snapshot-10 to staging!
```

Template | Content | Rendered example
---|---|---
`{{.Build.PWD}}`| Contents of the `PWD` environment variable. | `PWD=/go/src/github.com/joshdk/hub-comment`

#### Git Context

The `.Git` object contains git-specific values.

For example, the following template could be rendered as so:

```
Built branch {{.Git.Branch}} (@{{.Git.SHA}}).
```

```
Built branch feature/foo (@706b2cc37585547f286284a7b488f252c0cea482).
```

Template | Content | Rendered example
---|---|---
`{{.Git.Branch}}` | Current git branch. | `feature/foo`
`{{.Git.PR}}` | URL for the current GitHub PR. | `https://github.com/joshdk/hub-comment/pull/123`
`{{.Git.SHA}}` | Current git ref | `706b2cc37585547f286284a7b488f252c0cea482`
`{{.Git.Tag}}` | Current git tag. | `1.0.0`

#### Labels Context

The `.Build` object is an alphabetically sorted list of GitHub labels that are applied to the current PR.

The function `label` is defined, and returns `true` if the given label name is applied

```
{{if label "enhancement"}}

{{end}}
```

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[circleci-badge]:         https://circleci.com/gh/joshdk/hub-comment.svg?&style=shield
[circleci-link]:          https://circleci.com/gh/joshdk/hub-comment/tree/master
[github-downloads-badge]: https://img.shields.io/github/downloads/joshdk/hub-comment/total.svg
[github-release-badge]:   https://img.shields.io/github/release/joshdk/hub-comment/all.svg
[github-release-link]:    https://github.com/joshdk/hub-comment/releases/latest
[license-badge]:          https://img.shields.io/badge/license-MIT-green.svg
[license-file]:           https://github.com/joshdk/hub-comment/blob/master/LICENSE.txt
[license-link]:           https://opensource.org/licenses/MIT
