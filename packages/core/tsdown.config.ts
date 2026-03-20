import { defineConfig } from "tsdown";

export default defineConfig({
	// Entry point(s)
	entry: ["src/index.ts"],

	// Output formats
	format: ["esm", "cjs"],

	// Generate TypeScript declarations
	dts: true,

	// Generate sourcemaps
	sourcemap: true,

	// Clean dist folder before build
	clean: true,

	// Minify output (optional, disable for better debugging)
	minify: false,

	// Tree-shaking
	treeshake: true,

	// Target environment
	target: "es2022",

	// External dependencies (don't bundle)
	external: [],

	// Watch mode options
	watch: false,
});
