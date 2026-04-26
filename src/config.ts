import type { openai } from "@ai-sdk/openai"
import { Entry } from "@napi-rs/keyring"
import Conf from "conf"
import z from "zod"

const projectName = "howtf"

export const configSchema = z.object({
	provider: z.enum(["openai", "google", "anthropic"]),
	model: z.string<Parameters<typeof openai>[0]>(),
})

export const defaultConfig = {
	provider: "openai",
	model: "gpt-4.1",
} satisfies z.infer<typeof configSchema>

export const config = new Conf<z.infer<typeof configSchema>>({
	projectName,
	serialize: (value) => {
		return JSON.stringify(configSchema.parse(value), null, "\t")
	},
	deserialize: (value) => {
		return configSchema.parse(JSON.parse(value))
	},
})

export const saveApiKey = (profile: string, apiKey: string) => {
	const entry = new Entry(projectName, profile)
	entry.setPassword(apiKey)
}

export const loadApiKey = (profile: string) => {
	const entry = new Entry(projectName, profile)
	return entry.getPassword()
}

export const deleteApiKey = (profile: string) => {
	const entry = new Entry(projectName, profile)
	entry.deletePassword()
}
