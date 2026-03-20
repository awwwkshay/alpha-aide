import z from "zod/v4";

export const nodeEnvSchema = z
	.enum(["development", "testing", "production"])
	.default("production");
