import { outro } from "@clack/prompts"
import pc from "picocolors"

export const brandText = pc.bold(`HOW${pc.red("TF")}`)

export const exit = (code: number) => {
	outro(brandText)
	process.exit(code)
}
