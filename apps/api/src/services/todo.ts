import { todosTable } from "@awwwkshay/node-ts-core";
import { eq } from "drizzle-orm";

import { db } from "@/db";

type TodoSelect = typeof todosTable.$inferSelect;
type TodoInsert = typeof todosTable.$inferInsert;

export const createTodo = async (
	data: TodoInsert,
): Promise<TodoSelect | undefined> => {
	const [todo] = await db.insert(todosTable).values(data).returning();
	return todo;
};

export const getTodos = async (): Promise<TodoSelect[]> => {
	const todos = await db.select().from(todosTable);
	return todos;
};

export const getTodoById = async (
	id: string,
): Promise<TodoSelect | undefined> => {
	const [todo] = await db
		.select()
		.from(todosTable)
		.where(eq(todosTable.id, id));
	return todo;
};

export const updateTodo = async (
	id: string,
	data: Partial<TodoInsert>,
): Promise<TodoSelect | undefined> => {
	const [todo] = await db
		.update(todosTable)
		.set(data)
		.where(eq(todosTable.id, id))
		.returning();
	return todo;
};

export const deleteTodo = async (
	id: string,
): Promise<TodoSelect | undefined> => {
	const [todo] = await db
		.delete(todosTable)
		.where(eq(todosTable.id, id))
		.returning();
	return todo;
};
