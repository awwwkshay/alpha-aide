import { Hono } from "hono";

import {
	createTodo,
	getTodos,
	updateTodo,
	getTodoById,
	deleteTodo,
} from "@/services";

const todosRouter = new Hono();

todosRouter.post("/", async (c) => {
	const body = await c.req.json();
	const todo = await createTodo(body);
	c.status(201);
	return c.json({ message: "Created todo successfully", data: todo });
});

todosRouter.get("/", async (c) => {
	const todos = await getTodos();
	return c.json({ message: "Read list of todos", data: todos });
});

todosRouter.get("/:id", async (c) => {
	const { id } = c.req.param();
	const todo = await getTodoById(id);
	return c.json({ message: `Read todo with id: ${id}`, data: todo });
});

todosRouter.patch("/:id", async (c) => {
	const { id } = c.req.param();
	const body = await c.req.json();
	const todo = await updateTodo(id, body);
	return c.json({ message: `Updated todo with id: ${id}`, data: todo });
});

todosRouter.delete("/:id", async (c) => {
	const { id } = c.req.param();
	const todo = await deleteTodo(id);
	return c.json({ message: `Deleted todo with id: ${id}`, data: todo });
});

export { todosRouter };
