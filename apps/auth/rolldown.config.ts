import { defineConfig } from "rolldown";
import tsconfigPaths from "rollup-plugin-tsconfig-paths";

export default defineConfig({
	input: "./src/server.ts",
	platform: "node",
	plugins: [
		tsconfigPaths({
			tsConfigPath: "./tsconfig.json",
		}),
	],
	output: {
		dir: "dist",
		format: "esm",
		cleanDir: true,
	},
	watch: {
		include: "src/**",
		clearScreen: false,
		exclude: ["node_modules"],
	},
	// ⛑️ IMPORTANT FIX
	external: [
		"pg",
		"pg-pool",
		"pg-protocol",
		"pg-types",
		"better-auth",
		"drizzle-orm",
		"@hono/node-server",
		"hono",
	],
});
