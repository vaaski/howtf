#! /usr/bin/env node

import { config } from "./config"
import prompts from "prompts"

const main = async () => {
  let key = config.get("key")
  if (!key) {
    const response = await prompts({
      type: "text",
      name: "key",
      message: "Enter your OpenAI API key",
      validate: Boolean,
    })

    if (typeof response.key !== "string") throw new Error("No key provided")

    config.set("key", response.key)
    key = response.key
  }

  console.log("Your key is", key)
}

main()
export {}
