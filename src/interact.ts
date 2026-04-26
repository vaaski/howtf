import type { ModelMessage } from "ai"
import { box, isCancel, log, selectKey, spinner, text } from "@clack/prompts"
import clipboard from "clipboardy"
import pc from "picocolors"
import { x } from "tinyexec"
import { prompts } from "./prompts"
import { exit } from "./util"

export const acquireProblem = async () => {
	const problem = await text({
		message: "Describe your problem to generate a command.",
	})

	if (isCancel(problem)) {
		exit(0)
		return ""
	} else {
		return problem
	}
}

// --------------------------------------------------------------------------------------

export const generateSolution = async (problem: string) => {
	const spin = spinner()
	spin.start("generating command")

	const { result: { output }, history } = await prompts.generate(problem)
	const command = `${pc.bold(output.generated_command.binary)} ${output.generated_command.args.join(" ")}`
	const commandRaw = `${output.generated_command.binary} ${output.generated_command.args.join(" ")}`

	spin.stop(`${pc.gray(output.short_explanation)}`)

	box(command)

	return { command, commandRaw, output, history }
}

export const editSolution = async (existingHistory: ModelMessage[], edit: string) => {
	const spin = spinner()
	spin.start("generating edits")

	const { result: { output }, history } = await prompts.edit(existingHistory, edit)
	const command = `${pc.bold(output.generated_command.binary)} ${output.generated_command.args.join(" ")}`
	const commandRaw = `${output.generated_command.binary} ${output.generated_command.args.join(" ")}`

	spin.stop(`${pc.gray(output.short_explanation)}`)

	box(command)

	return { command, commandRaw, output, history }
}

// --------------------------------------------------------------------------------------

type Solution = Awaited<ReturnType<typeof generateSolution>>
export const useSolution = async (
	{ command, commandRaw, output, history }: Solution,
) => {
	const action = await selectKey({
		message: "What next?",
		options: [
			{ value: "r", label: "Run" },
			{ value: "e", label: "Edit" },
			{ value: "c", label: "Copy" },
			{ value: "x", label: "Exit" },
		],
		caseSensitive: false,
	})

	if (isCancel(action)) {
		exit(0)
	} else if (action === "r") {
		log.info(`running ${command}`)

		const child = x(output.generated_command.binary, output.generated_command.args, {
			nodeOptions: {
				cwd: process.cwd(),
				env: process.env,
				stdio: ["ignore", "inherit", "inherit"],
			},
		})

		await child
	} else if (action === "e") {
		const editPrompt = await text({
			message: "What needs changing?.",
		})

		if (isCancel(editPrompt)) {
			exit(0)
		} else {
			await useSolution(await editSolution(history, editPrompt))
		}
	} else if (action === "c") {
		await clipboard.write(commandRaw)
		log.success(`copied ${command}`)
	}
}
