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

export const generateSolution = async (problem: string) => {
	const spin = spinner()
	spin.start("generating command")

	const { output } = await prompts.generate(problem)
	const command = `${pc.bold(output.generated_command.binary)} ${output.generated_command.args.join(" ")}`
	const commandRaw = `${output.generated_command.binary} ${output.generated_command.args.join(" ")}`

	spin.stop(`${pc.gray(output.short_explanation)}`)

	box(command)

	return { command, commandRaw, output }
}

export const editSolution = async (problem: string, solution: string, edit: string) => {
	const spin = spinner()
	spin.start("generating edits")

	const { output } = await prompts.edit(problem, solution, edit)
	const command = `${pc.bold(output.generated_command.binary)} ${output.generated_command.args.join(" ")}`
	const commandRaw = `${output.generated_command.binary} ${output.generated_command.args.join(" ")}`

	spin.stop(`${pc.gray(output.short_explanation)}`)

	box(command)

	return { command, commandRaw, output }
}

export const useSolution = async (problem: string, {
	command,
	commandRaw,
	output,
}: Awaited<ReturnType<typeof generateSolution>>) => {
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
		const edit = await text({
			message: "What needs changing?.",
		})

		if (isCancel(edit)) {
			exit(0)
		} else {
			await useSolution(problem, await editSolution(problem, command, edit))
		}
	} else if (action === "c") {
		await clipboard.write(commandRaw)
		log.success(`copied ${command}`)
	}
}
