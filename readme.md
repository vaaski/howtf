<h1 align="center">howtf</h1>
<p align="center">quickly look up console commands using the GPT API</p>

<p align="center">
  <img src="docs/ffmpeg-example.svg" width="600" align="center" />
</p>

---

### Synopsis

This is a pretty simple CLI to help you find that command you can never
remember using the OpenAI chat completion API with either gpt-3.5-turbo or gpt-4.
Usually faster than googling it.

Also comes with the neat advantage of being personalized to whatever you need.
Notice how in the example above, I specifically asked to modify "Video.mov" and
it automatically generated Video_stretched.mov as the output file.

### Installation

```bash
npm i -g howtf
```

Save your OpenAI API key:

```bash
howtf -K [your API key]
```

Or just use it once instead of saving it it:

```bash
howtf -k [your API key] [command]
```

### Usage

```bash
howtf [command]
```

For example:

```bash
howtf convert input.mp4 to gif with ffmpeg
# ❯ ffmpeg -i input.mp4 output.gif
```

```bash
howtf undo last 3 commits
# ❯ git reset HEAD~3 --hard
```