import Conf from "conf"

export type HowtfRequest = {
  query: string
  response: string
}

export type Model = "gpt-3.5-turbo" | "gpt-4"

const defaults = {
  key: undefined as string | undefined,
  model: "gpt-3.5-turbo" as Model,
  lastRequest: undefined as HowtfRequest | undefined,
}

export const config = new Conf({
  projectName: "howtf",
  defaults,
})
