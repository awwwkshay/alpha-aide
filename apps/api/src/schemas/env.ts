import z from "zod/v4";

export const processEnvSchema = z.object({
	NODE_ENV: z
		.enum(["development", "testing", "production"])
		.default("production"),
	PORT: z.coerce.number().default(8000),
	CLIENT_URLS: z
		.string()
		.transform((val) => val.split(","))
		.default([]),
	DATABASE_URL: z.url(),
	DB_SSL: z.coerce.boolean().default(false),
});

export type ProcessEnv = z.infer<typeof processEnvSchema>;
