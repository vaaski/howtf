await Bun.build({
	entrypoints: ["./index.ts"],
	outdir: "./dist",
	target: "node",
	external: ["@napi-rs/keyring"],
})
