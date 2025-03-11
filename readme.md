<h1 align="center">howtf</h1>
<p align="center">Quickly look up console commands using the OpenAI GPT API</p>

<p align="center">
  <img src="docs/ffmpeg-example.svg" width="600" align="center" />
</p>

---

### Synopsis

This is a pretty simple CLI tool to help you find that command you can never
remember using the OpenAI chat completion API. This usually ends up being
faster than googling it.

Also comes with the neat advantage of being personalized to whatever you need.
Notice how in the example above, I specifically asked to modify "example.webm" and
it automatically chose "example.mp4" as the output file.

### Features

- 🔒 Securely stores your OpenAI API key in your OS keychain with [go-keyring](https://github.com/zalando/go-keyring)
- 🙈 Lets you inspect the generated command before executing or copying it
- 💪 Lets you describe your command as arguments or interactively to avoid escaping special characters
- ⁉️ Can explain what any command does
- 💅 Comes with a beautiful Terminal User Interface
- 🆖 Cross-platform compatible

### Installation

Install using `go install` (requires [Go](https://go.dev/doc/install)):

```bash
go install github.com/vaaski/howtf@go
```

Run the interactive configuration wizard to set your OpenAI API key:

```bash
howtf -config
```

### Usage

```bash
howtf [command]
```

### Examples

Interactive mode:

<img src="docs/interactive.svg" width="600" />

<br>

Config page:

<img src="docs/config.svg" width="600" />
