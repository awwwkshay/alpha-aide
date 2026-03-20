import { ScriptOnce } from "@tanstack/react-router";
import { createClientOnlyFn, createIsomorphicFn } from "@tanstack/react-start";
import { createContext, ReactNode, useEffect, useState } from "react";

import { AppTheme, UserTheme, userThemeSchema } from "@/schemas";

const themeStorageKey = "ui-theme";

const getStoredUserTheme = createIsomorphicFn()
	.server((): UserTheme => "system")
	.client((): UserTheme => {
		const stored = localStorage.getItem(themeStorageKey);
		if (!stored) return "system";

		// Handle cases where the stored value might be a JSON array (e.g. ["light"])
		// or a plain string ("light"). Normalize to a string before validation.
		let value: unknown = stored;
		try {
			const parsed = JSON.parse(stored);
			if (Array.isArray(parsed)) {
				value = parsed[0];
			} else if (typeof parsed === "string") {
				value = parsed;
			} else {
				value = stored;
			}
		} catch {
			// If JSON.parse fails, assume the stored value is a plain string
			value = stored;
		}

		return userThemeSchema.parse(value as any);
	});

const setStoredTheme = createClientOnlyFn((theme: UserTheme) => {
	const validatedTheme = userThemeSchema.parse(theme);
	localStorage.setItem(themeStorageKey, validatedTheme);
});

const getSystemTheme = createIsomorphicFn()
	.server((): AppTheme => "light")
	.client((): AppTheme => {
		return window.matchMedia("(prefers-color-scheme: dark)").matches
			? "dark"
			: "light";
	});

const handleThemeChange = createClientOnlyFn((userTheme: UserTheme) => {
	const validatedTheme = userThemeSchema.parse(userTheme);

	const root = document.documentElement;
	root.classList.remove("light", "dark", "system");

	if (validatedTheme === "system") {
		const systemTheme = getSystemTheme();
		root.classList.add(systemTheme, "system");
	} else {
		root.classList.add(validatedTheme);
	}
});

const setupPreferredListener = createClientOnlyFn(() => {
	const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
	const handler = () => handleThemeChange("system");
	mediaQuery.addEventListener("change", handler);
	return () => mediaQuery.removeEventListener("change", handler);
});

const themeScript = (function () {
	function themeFn() {
		try {
			const storedTheme = localStorage.getItem("ui-theme") || "system";
			const validTheme = ["light", "dark", "system"].includes(storedTheme)
				? storedTheme
				: "system";

			if (validTheme === "system") {
				const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
					.matches
					? "dark"
					: "light";
				document.documentElement.classList.add(systemTheme, "system");
			} else {
				document.documentElement.classList.add(validTheme);
			}
		} catch {
			const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
				.matches
				? "dark"
				: "light";
			document.documentElement.classList.add(systemTheme, "system");
		}
	}
	return `(${themeFn.toString()})();`;
})();

type ThemeContextProps = {
	userTheme: UserTheme;
	appTheme: AppTheme;
	setTheme: (theme: UserTheme) => void;
};
export const ThemeContext = createContext<ThemeContextProps | undefined>(
	undefined,
);

type ThemeProviderProps = {
	children: ReactNode;
};
export function ThemeProvider({ children }: ThemeProviderProps) {
	const [userTheme, setUserTheme] = useState<UserTheme>(getStoredUserTheme);

	useEffect(() => {
		if (userTheme !== "system") return;
		return setupPreferredListener();
	}, [userTheme]);

	const appTheme = userTheme === "system" ? getSystemTheme() : userTheme;

	const setTheme = (newUserTheme: UserTheme) => {
		const validatedTheme = userThemeSchema.parse(newUserTheme);
		setUserTheme(validatedTheme);
		setStoredTheme(validatedTheme);
		handleThemeChange(validatedTheme);
	};

	return (
		<ThemeContext value={{ userTheme, appTheme, setTheme }}>
			<ScriptOnce children={themeScript} />

			{children}
		</ThemeContext>
	);
}
