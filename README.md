[![CircleCI][circleci-badge]][circleci-link]
[![License][license-badge]][license-link]
[![Github downloads][github-downloads-badge]][github-release-link]
[![GitHub release][github-release-badge]][github-release-link]

# Hub-Comment

📝 Automate GitHub pull request comments in CI

## Installing

### From source

You can install a development version of this tool by running:

```bash
$ go get -u github.com/joshdk/hub-comment
```

### Precompiled binary

Alternatively, you can download a precompiled [release][github-release-link] binary by running:

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

Build #{{.Env.CIRCLE_BUILD_NUM}} was run successfully.
```

### Running

```
$ hub-comment -template-file hello-template.txt
Posting new comment as Josh Komoroske (joshdk) on joshdk/hub-comment#123:

→ 🤖 Hello from CI!
→
→ Build #123 was run successfully.

To view comment visit:
→ https://github.com/joshdk/hub-comment/issues/123#issuecomment-421483151
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
