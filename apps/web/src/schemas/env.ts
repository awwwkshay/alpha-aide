// src/config/env.ts
import { z } from "zod/v4";

export const envSchema = z.object({
	NITRO_API_BASE_URL: z.string(),
	NITRO_PORT: z.string().transform((port) => parseInt(port)),
});

export const clientEnvSchema = z.object({
	VITE_APP_NAME: z.string(),
	VITE_API_BASE_URL: z.url(),
});
