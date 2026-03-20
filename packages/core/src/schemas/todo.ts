import {
	createSelectSchema,
	createUpdateSchema,
	createInsertSchema,
} from "drizzle-zod";
import z from "zod/v4";

import { todosTable } from "../db/schema";

const title = z.string().min(3).max(255);

export const todoSchema = createSelectSchema(todosTable).extend({
	title,
});

export type ITodo = z.infer<typeof todoSchema>;

export const todoUpdateSchema = createUpdateSchema(todosTable).extend({
	title: title.optional(),
});

export type ITodoUpdate = z.infer<typeof todoUpdateSchema>;

export const todoInsertSchema = createInsertSchema(todosTable).extend({
	title,
});

export type ITodoCreate = z.infer<typeof todoInsertSchema>;
