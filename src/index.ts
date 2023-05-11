#! /usr/bin/env node

import type { Flags } from "./config"

import { execaCommand } from "execa"
import minimist from "minimist"
import ora from "ora"
import prompts from "prompts"
import { z } from "zod"

import packageJson from "../package.json"
import { config, humanModelEnum } from "./config"
import { askAI } from "./gpt"
import { modelShortcuts, printHelp } from "./util"

const parameters = minimist<Flags>(process.argv.slice(2), {
  string: ["key", "k", "model", "m", "set-model", "M", "set-key", "K"],
})
const query = parameters._.join(" ")

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

  const modelParameter =
    parameters.model || parameters.m || parameters["set-model"] || parameters.M
  const saveModel = Boolean(parameters["set-model"] || parameters.M)
  const suppliedModel = humanModelEnum.safeParse(modelParameter)

  if (suppliedModel.success) {
    model = modelShortcuts(suppliedModel.data)
    if (saveModel) {
      config.set("model", model)
      return console.log(`Default model saved as ${model}`)
    }
  } else if (modelParameter) {
    throw new Error(`Invalid model: ${modelParameter}`)
  }

  if (parameters.help || parameters.h || parameters._.length === 0) {
    return printHelp()
  }

  if (!key) {
    const keyResponse = await promptKey()
    key = z.string().parse(keyResponse.key)
  }

  const spinner = ora("Generating command...").start()
  const command = await askAI(query, model, key)
  spinner.succeed(`generated command:\n\n  ❯ ${command}\n`)

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

    await execaCommand(command, { stdio: "inherit", shell: true })
  }
}

main()
