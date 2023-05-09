#! /usr/bin/env node

import minimist from "minimist"
import prompts from "prompts"
import { z } from "zod"
import { execaCommand } from "execa"
import packageJson from "../package.json"
import { config, modelEnum } from "./config"
import { askAI } from "./gpt"

const parameters = minimist(process.argv.slice(2), { string: ["key", "k", "model", "m"] })
const query = parameters._.join(" ")

const help = () => {
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
    "  print contents of all files ending in .js in current directory",
    "",
    "Config file location:",
    `  ${config.path}`,
  ].join("\n")

  console.log(helpText)
}

const promptKey = async () => {
  const response = await prompts({
    type: "password",
    name: "key",
    message: "Enter your OpenAI API key",
    validate: Boolean,
  })

  if (typeof response.key !== "string") throw new Error("No key provided")
  return response
}

const main = async () => {
  let key = config.get("key")
  let model = config.get("model")

  if (parameters.version || parameters.v) {
    return console.log(packageJson.version)
  }
  if (parameters.config || parameters.c) {
    return console.log(config.path)
  }
  if (parameters.clear || parameters.C) {
    return config.clear()
  }

  const suppliedKey = parameters.key || parameters.k
  if (suppliedKey && typeof suppliedKey === "string") {
    key = suppliedKey
  }

  const suppliedSetKey = parameters["set-key"] || parameters.K
  if (suppliedSetKey) {
    if (typeof suppliedSetKey === "string") {
      config.set("key", suppliedSetKey)
      key = suppliedSetKey
    } else {
      const keyResponse = await promptKey()
      config.set("key", keyResponse.key)
      key = keyResponse.key
    }

    return console.log("Key saved")
  }

  if (parameters.help || parameters.h || parameters._.length === 0) {
    return help()
  }

  if (!key) {
    const keyResponse = await promptKey()
    key = z.string().parse(keyResponse.key)
  }

  const modelParameter = parameters.model || parameters.m
  const suppliedModel = modelEnum.safeParse(modelParameter)
  if (suppliedModel.success) {
    if (suppliedModel.data === "3") model = "gpt-3.5-turbo"
    else if (suppliedModel.data === "4") model = "gpt-4"
    else model = suppliedModel.data
  } else if (modelParameter) {
    throw new Error(`Invalid model: ${modelParameter}`)
  }

  const command = await askAI(query, model, key)
  console.log(`generated command:\n\n  ❯ ${command}\n`)

  const shouldExecute = await prompts({
    name: "execute",
    type: "confirm",
    message: "Execute it?",
    initial: true,
  })

  console.log(shouldExecute.execute ? "executing..." : "not executing")

  if (shouldExecute.execute) {
    process.stdout.moveCursor(0, -2)
    process.stdout.clearScreenDown()

    await execaCommand(command, { stdio: "inherit" })
  }
}

main()
