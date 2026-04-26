import { box, intro, isCancel, log, outro, selectKey, spinner } from "@clack/prompts"
import clipboard from "clipboardy"
import pc from "picocolors"
import { x } from "tinyexec"
import { prompts } from "./src/prompts"
import { introText } from "./src/util"
import { configView } from "./views/config"

const args = process.argv.slice(2)

intro(introText)

switch (args.join("")) {
	case "config":
		await configView()
		break

	case "":
		console.log("interactive mode")
		break

	default: {
		const problem = args.join(" ")

		const spin = spinner()
		spin.start("generating command")

		const { output } = await prompts.generate(problem)
		const command = `${pc.bold(output.generated_command.binary)} ${output.generated_command.args.join(" ")}`

		spin.stop(`${pc.gray(output.short_explanation)}`)

		box(command)

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
			outro(introText)
			process.exit(0)
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
			log.error(`edit ${command}`)
		} else if (action === "c") {
			await clipboard.write(command)
			log.success(`copied ${command}`)
		}

		break
	}
}

outro(introText)
