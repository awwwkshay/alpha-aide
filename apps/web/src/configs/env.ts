import { clientEnvSchema, envSchema } from "@/schemas";

// Validate server environment
export const serverEnv = envSchema.parse(process.env || {});

// Validate client environment
export const clientEnv = clientEnvSchema.parse(import.meta.env || {});
