import { pgTable, varchar, uuid, boolean } from "drizzle-orm/pg-core";

export const todosTable = pgTable("todos", {
	id: uuid().primaryKey().defaultRandom(),
	title: varchar({ length: 255 }).notNull(),
	completed: boolean().default(false),
});
