import { platform } from "node:os"
import { createOpenAI } from "@ai-sdk/openai"
import { generateText, Output } from "ai"
import { z } from "zod"
import { loadApiKey } from "./config"

const openai = createOpenAI({
	apiKey: loadApiKey("openai") ?? "",
})

export const prompts = {
	generate: (problem: string) => {
		return generateText({
			model: openai("gpt-4.1"),
			output: Output.object({
				schema: z.object({
					short_explanation: z.string(),
					generated_command: z.object({
						binary: z.string(),
						args: z.array(z.string()),
					}),
				}),
			}),
			messages: [
				{
					role: "system",
					content: [
						"You generate oneliner-commands for a given problem.",
						`The user is running ${process.env.SHELL} on ${platform()}`,
					].join("\n"),
				},
				{
					role: "user",
					content: problem,
				},
			],
		})
	},
	edit: (problem: string, solution: string, edit: string) => {
		return generateText({
			model: openai("gpt-4.1"),
			output: Output.object({
				schema: z.object({
					short_explanation: z.string(),
					generated_command: z.object({
						binary: z.string(),
						args: z.array(z.string()),
					}),
				}),
			}),
			messages: [
				{
					role: "system",
					content: [
						"You generate oneliner-commands for a given problem.",
						`The user is running ${process.env.SHELL} on ${platform()}`,
					].join("\n"),
				},
				{
					role: "user",
					content: problem,
				},
				{
					role: "assistant",
					content: solution,
				},
				{
					role: "user",
					content: edit,
				},
			],
		})
	},
}
