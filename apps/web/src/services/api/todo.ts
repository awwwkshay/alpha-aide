import {
	ITodo,
	ITodoCreate,
	ITodoUpdate,
	IApiResponse,
	todoSchema,
} from "@awwwkshay/node-ts-core";
import axios from "axios";
import z from "zod/v4";

// Construct schemas
const todoArraySchema = z.array(todoSchema);

// Construct TS types from schemas
type TodoResponse = IApiResponse<typeof todoSchema>;
type TodoArrayResponse = IApiResponse<typeof todoArraySchema>;

// ---- API FUNCTIONS ----

export const createTodo = async (
	todo: ITodoCreate,
): Promise<ITodo | undefined> => {
	const res = await axios.post<TodoResponse>("/todos", todo);
	return res.data.data;
};

export const readTodos = async (): Promise<ITodo[] | undefined> => {
	const res = await axios.get<TodoArrayResponse>("/todos");
	return res.data.data;
};

export const readTodoById = async (id: string): Promise<ITodo | undefined> => {
	const res = await axios.get<TodoResponse>(`/todos/${id}`);
	return res.data.data;
};

export const updateTodo = async (
	id: string,
	updates: ITodoUpdate,
): Promise<ITodo | undefined> => {
	const res = await axios.patch<TodoResponse>(`/todos/${id}`, updates);
	return res.data.data;
};

export const deleteTodo = async (id: string): Promise<ITodo | undefined> => {
	const res = await axios.delete<TodoResponse>(`/todos/${id}`);
	return res.data.data;
};
