import tailwindcss from "@tailwindcss/vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { nitro } from "nitro/vite";
import { defineConfig } from "vite";
import viteTsConfigPaths from "vite-tsconfig-paths";

const config = defineConfig({
	plugins: [
		devtools(),
		nitro({
			devServer: {
				hostname: "0.0.0.0",
				port: 3000,
			},
			output: {
				dir: "./dist",
				publicDir: "./dist/public",
				serverDir: "./dist/server",
			},
		}),
		// this is the plugin that enables path aliases
		viteTsConfigPaths({
			projects: ["./tsconfig.json"],
		}),
		tailwindcss(),
		tanstackStart(),
		viteReact(),
	],

	server: {
		host: true,
		port: 3000,
		hmr: {
			host: "0.0.0.0",
		},
	},
});

export default config;
