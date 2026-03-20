import { processEnvSchema } from "./schemas";

if (process.env.NODE_ENV === "development") {
	console.log("Loading .env file for development environment");
	await import("dotenv/config");
}

export const ENV_VARS = processEnvSchema.parse(process.env);
