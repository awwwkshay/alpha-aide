import z from "zod/v4";

const light = z.literal("light");
const dark = z.literal("dark");
const system = z.literal("system");

export const userThemeSchema = z.union([light, dark, system]);

export type UserTheme = z.infer<typeof userThemeSchema>;

export const appThemeSchema = z.union([light, dark]);

export type AppTheme = z.infer<typeof appThemeSchema>;
