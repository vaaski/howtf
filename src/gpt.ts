import { Configuration, OpenAIApi } from "openai"

let OS_HUMAN: string = process.platform
if (OS_HUMAN === "darwin") OS_HUMAN = "macos"
if (OS_HUMAN === "win32") OS_HUMAN = "windows"

const SYSTEM_MESSAGE = `You're a generator for console commands for a developer. Do not explain anything, just provide the command.`
const PRETEXT = `I want to solve the following problem in the console on ${OS_HUMAN}:\n`

const markdownRegex = /```(?:.+\n)?([\S\s]+)\n```|`(.+)`/m
const extractMarkdownMaybe = (text: string) => {
  const match = text.match(markdownRegex)
  if (!match) return text

  return match[1] || match[2]
}

export const askAI = async (query: string, model: string, apiKey: string) => {
  const configuration = new Configuration({ apiKey })
  const openai = new OpenAIApi(configuration)

  const response = await openai.createChatCompletion({
    model,
    temperature: 0,
    max_tokens: 128,
    messages: [
      { role: "system", content: SYSTEM_MESSAGE },
      { role: "user", content: PRETEXT + query },
    ],
  })

  // console.log(JSON.stringify(response.data, undefined, 2))
  // console.log(response.data.choices[0].message?.content)

  const answer = response.data.choices[0]
  const command = answer.message?.content
  if (!command) throw new Error("No command generated")

  return extractMarkdownMaybe(command)
}
