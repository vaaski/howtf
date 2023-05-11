import type { HumanModel, Model } from "./config"
import { config } from "./config"

export const modelShortcuts = (model: HumanModel): Model => {
  if (model === "3") return "gpt-3.5-turbo"
  if (model === "4") return "gpt-4"

  return model
}

export const printHelp = () => {
  const helpText = [
    "howtf - use GPT to quickly look up a console command",
    "",
    "Usage:",
    "  howtf [command]",
    "",
    "Options:",
    "  --help, -h       print help",
    "  --version, -v    print version",
    "  --model, -m      use model (gpt-3.5-turbo or gpt-4, you can also just pass 3 or 4)",
    "  --set-model, -M  set model and save to config",
    "  --key, -k        pass OpenAI API key once",
    "  --set-key, -K    set OpenAI API key",
    "  --config, -c     print config file location",
    "  --clear, -C      clear config file",
    "",
    "Examples:",
    "  howtf undo last 2 git commits",
    "  howtf remove all untracked files in git",
    "  howtf print contents of all files ending in .js in current directory",
    "",
    "Config file location:",
    `  ${config.path}`,
  ].join("\n")

  console.log(helpText)
}
