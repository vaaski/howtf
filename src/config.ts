import Conf from "conf"
import { z } from "zod"

export type HowtfRequest = {
  query: string
  response: string
}

export const modelEnum = z.enum(["gpt-3.5-turbo", "gpt-4"])
export type Model = z.infer<typeof modelEnum>

export const humanModelEnum = modelEnum.or(z.enum(["3", "4"]))
export type HumanModel = z.infer<typeof humanModelEnum>

const defaults = {
  key: undefined as string | undefined,
  model: "gpt-3.5-turbo" as Model,
  lastRequest: undefined as HowtfRequest | undefined,
}

export const config = new Conf({
  projectName: "howtf",
  defaults,
})

// todo: zoddify this
export type Flags = {
  h?: boolean
  help?: boolean
  v?: boolean
  version?: boolean
  m?: string
  model?: string
  M?: string
  "set-model"?: string
  k?: string
  key?: string
  K?: string
  "set-key"?: string
  c?: boolean
  config?: boolean
  C?: boolean
  clear?: boolean
}
