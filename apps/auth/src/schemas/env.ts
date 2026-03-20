import z from "zod/v4";

export const processEnvSchema = z.object({
	NODE_ENV: z
		.enum(["development", "testing", "production"])
		.default("production"),
	PORT: z.coerce.number().default(4000),
	CLIENT_URLS: z
		.string()
		.transform((val) => val.split(","))
		.default([]),
	BETTER_AUTH_SECRET: z.string().min(18),
	BETTER_AUTH_URL: z.url(),
	DATABASE_URL: z.url(),
});

export type ProcessEnv = z.infer<typeof processEnvSchema>;
