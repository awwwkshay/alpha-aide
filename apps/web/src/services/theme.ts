import { createServerFn } from "@tanstack/react-start";
import { getCookie, setCookie } from "@tanstack/react-start/server";

import { UserTheme, userThemeSchema } from "@/schemas";

const storageKey = "_preferred-theme";

export const getThemeServerFn = createServerFn().handler(
	async () => (getCookie(storageKey) || "system") as UserTheme,
);

export const setThemeServerFn = createServerFn({ method: "POST" })
	.inputValidator(userThemeSchema)
	.handler(async ({ data }) => setCookie(storageKey, data));
