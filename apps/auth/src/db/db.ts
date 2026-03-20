import { drizzle } from "drizzle-orm/node-postgres";
import { Pool } from "pg";

import { ENV_VARS } from "@/config";
import { accounts, sessions, users, verifications, jwkss } from "@/db/schema";

export const pool = new Pool({
	connectionString: ENV_VARS.DATABASE_URL,
	max: 10,
	idleTimeoutMillis: 30000,
});

export const db = drizzle(pool, {
	schema: {
		accounts,
		sessions,
		users,
		verifications,
		jwkss,
	},
	casing: "snake_case",
});
