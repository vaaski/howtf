import type { ModelMessage } from "ai"
import { platform } from "node:os"
import { createOpenAI } from "@ai-sdk/openai"
import { generateText, Output } from "ai"
import { z } from "zod"
import { loadApiKey } from "./config"

const openai = createOpenAI({
	apiKey: loadApiKey("openai") ?? "",
})

export const prompts = {
	generate: async (problem: string) => {
		const messages: ModelMessage[] = [
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
		]

		const result = await generateText({
			messages,
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
		})

		return {
			result,
			history: [
				...messages,
				...result.response.messages,
			],
		}
	},
	edit: async (history: ModelMessage[], edit: string) => {
		const messages: ModelMessage[] = [
			...history,
			{
				role: "user",
				content: edit,
			},
		]

		const result = await generateText({
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
			messages,
		})

		return {
			result,
			history: [
				...messages,
				...result.response.messages,
			],
		}
	},
}
