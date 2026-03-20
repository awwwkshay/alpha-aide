import { drizzle } from "drizzle-orm/node-postgres";

import { ENV_VARS } from "@/config";

import { todosTable } from "./schema";

const useSsl = ENV_VARS.DB_SSL;

const connectionOptions = {
	connectionString: ENV_VARS.DATABASE_URL,
	ssl: useSsl ? { rejectUnauthorized: false } : false,
};

export const db = drizzle({
	connection: connectionOptions,
	schema: {
		todos: todosTable,
	},
});
