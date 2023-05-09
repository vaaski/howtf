import Conf from "conf"
import { z } from "zod"

export type HowtfRequest = {
  query: string
  response: string
}

export const modelEnum = z.enum(["gpt-3.5-turbo", "gpt-4", "3", "4"])
export type Model = z.infer<typeof modelEnum>

const defaults = {
  key: undefined as string | undefined,
  model: "gpt-3.5-turbo" as Model,
  lastRequest: undefined as HowtfRequest | undefined,
}

export const config = new Conf({
  projectName: "howtf",
  defaults,
})
