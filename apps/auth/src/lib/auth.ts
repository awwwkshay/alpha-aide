import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { jwt, openAPI } from "better-auth/plugins";

import { ENV_VARS } from "@/config";
import { accounts, sessions, users, verifications, db, jwkss } from "@/db";

export const auth = betterAuth({
	trustedOrigins: ["*"],
	secret: ENV_VARS.BETTER_AUTH_SECRET,
	url: ENV_VARS.BETTER_AUTH_URL,
	basePath: "/",
	user: {
		deleteUser: {
			enabled: true,
		},
	},
	database: drizzleAdapter(db, {
		provider: "pg",
		usePlural: true,
		schema: {
			accounts,
			sessions,
			users,
			verifications,
			jwkss,
		},
	}),
	appName: "auth",
	plugins: [
		openAPI({
			theme: "deepSpace",
			path: "docs",
		}),
		jwt(),
	],
	emailAndPassword: {
		enabled: true,
	},
	advanced: {
		cookies: {
			sessionToken: {
				attributes: {
					sameSite: "none",
					secure: true,
					partitioned: true, // New browser standards will mandate this for foreign cookies
				},
			},
		},
	},
});
