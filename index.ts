import { intro } from "@clack/prompts"
import { acquireProblem, generateSolution, useSolution } from "./src/interact"
import { brandText, exit } from "./src/util"
import { configView } from "./views/config"

const args = process.argv.slice(2)

intro(brandText)

switch (args.join("")) {
	case "config":
		await configView()
		break

	case "": {
		const problem = await acquireProblem()
		const solution = await generateSolution(problem)
		await useSolution(problem, solution)
		break
	}

	default: {
		const problem = args.join(" ")
		const solution = await generateSolution(problem)
		await useSolution(problem, solution)

		break
	}
}

exit(0)
