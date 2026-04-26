import { isCancel, log, note, password, selectKey } from "@clack/prompts"
import pc from "picocolors"
import { deleteApiKey, loadApiKey, saveApiKey } from "../src/config"

export const configView = async () => {
	const existingKey = loadApiKey("openai")

	note(`OpenAI API key: ${existingKey ?? pc.gray("not set")}`, "current config")

	const action = await selectKey({
		message: "What next?",
		options: [
			{ value: "e", label: "Edit API key" },
			{ value: "d", label: "Delete API key" },
			{ value: "x", label: "Exit" },
		],
		caseSensitive: false,
	})

	if (isCancel(action) || action === "x") {
		process.exit(0)
	}

	if (action === "d") {
		deleteApiKey("openai")
		log.success("API key deleted")

		process.exit(0)
	}

	const apiKey = await password({
		message: "OpenAI API key",
	})

	if (isCancel(apiKey)) {
		process.exit(0)
	}

	if (apiKey) {
		saveApiKey("openai", apiKey)
		log.success("API key saved")
	}
}
