import { ITodo, todoInsertSchema } from "@awwwkshay/node-ts-core";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useServerFn } from "@tanstack/react-start";
import { Pencil, Trash2, Check, X, CheckCheck, Undo2 } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import { useState } from "react";

import { useAppForm } from "@/lib";
import { cn } from "@/lib/utils";
import { createTodo, deleteTodo, readTodos, updateTodo } from "@/services";

import { Button } from "./ui/button";
import { Input } from "./ui/input";

export const TodoPanel = () => {
	const createTodoFn = useServerFn(createTodo);
	const updateTodoFn = useServerFn(updateTodo);
	const deleteTodoFn = useServerFn(deleteTodo);
	const readTodosFn = useServerFn(readTodos);

	const queryClient = useQueryClient();
	const [editingId, setEditingId] = useState<string | null>(null);
	const [editValue, setEditValue] = useState("");

	// READ TODOS
	const todosQuery = useQuery({
		queryKey: ["todos"],
		queryFn: readTodosFn,
	});

	// MUTATIONS
	const createTodoMutation = useMutation({
		mutationFn: (data: { title: string }) => createTodoFn(data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
	});

	const updateTodoMutation = useMutation({
		mutationFn: (data: { id: string; updates: Partial<ITodo> }) =>
			updateTodoFn(data.id, data.updates),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
	});

	const deleteTodoMutation = useMutation({
		mutationFn: (id: string) => deleteTodoFn(id),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
	});

	// CREATE FORM
	const createTodoForm = useAppForm({
		defaultValues: { title: "" },
		validators: { onChange: todoInsertSchema },
		onSubmit: async ({ value, formApi }) => {
			await createTodoMutation.mutateAsync({ title: value.title });
			formApi.reset();
			// Reset touched so button is enabled again for empty input
			formApi.setFieldMeta("title", (meta) => ({ ...meta, touched: false }));
		},
	});

	const { Field, handleSubmit } = createTodoForm;

	// Animation variants
	const itemVariants = {
		hidden: { opacity: 0, y: -10 },
		visible: { opacity: 1, y: 0 },
		exit: { opacity: 0, y: 10, transition: { duration: 0.2 } },
	};

	return (
		<div className="w-full flex flex-col gap-4 p-2">
			<h2 className="text-xl font-semibold">Todos</h2>

			{/* CREATE TODO */}
			<div className="flex flex-row gap-2 items-start">
				<div className="flex-1">
					<Field name="title">
						{(field) => {
							const meta = field.state.meta;
							const errorMessage = meta.errors?.[0]?.message;

							return (
								<div className="flex flex-col gap-1 w-full">
									<Input
										value={field.state.value}
										placeholder="New todo..."
										onChange={(e) => field.handleChange(e.target.value)}
										onBlur={() => field.handleBlur()}
										className="border rounded w-full focus-visible:ring-1"
									/>
									{meta.isTouched && errorMessage && (
										<p className="text-xs text-red-500 px-1">{errorMessage}</p>
									)}
								</div>
							);
						}}
					</Field>
				</div>

				{/* ADD BUTTON LOGIC: disable only when touched + dirty + invalid */}
				<Field name="title">
					{(field) => {
						const meta = field.state.meta;
						const disabled =
							meta.isTouched && meta.isDirty && meta.errors.length > 0;

						return (
							<Button onClick={() => handleSubmit()} disabled={disabled}>
								Add
							</Button>
						);
					}}
				</Field>
			</div>

			{/* TODOS LIST */}
			<div className="flex flex-col gap-2">
				{todosQuery.isLoading && <p>Loading...</p>}

				<AnimatePresence>
					{todosQuery.data?.map((todo) => {
						const isEditing = editingId === todo.id;
						const isChanged = editValue.trim() !== todo.title.trim();

						return (
							<motion.div
								key={todo.id}
								variants={itemVariants}
								initial="hidden"
								animate="visible"
								exit="exit"
								layout
								className={cn(
									"p-2 border rounded flex flex-row items-center gap-2",
									todo.completed && "line-through opacity-50",
								)}
							>
								{/* TITLE / EDIT MODE */}
								<div className="flex-1">
									{isEditing ? (
										<Input
											value={editValue}
											onChange={(e) => setEditValue(e.target.value)}
											className="w-full border rounded focus-visible:ring-1"
										/>
									) : (
										<span>{todo.title}</span>
									)}
								</div>

								{/* ACTION BUTTONS */}
								<div className="flex flex-row gap-2">
									{!isEditing && (
										<>
											<Button
												variant="destructive"
												size="sm"
												onClick={() => deleteTodoMutation.mutate(todo.id)}
											>
												<Trash2 className="h-4 w-4" />
											</Button>

											<Button
												variant="secondary"
												size="sm"
												onClick={() => {
													setEditValue(todo.title);
													setEditingId(todo.id);
												}}
											>
												<Pencil className="h-4 w-4" />
											</Button>

											{todo.completed ? (
												<Button
													size="sm"
													onClick={() =>
														updateTodoMutation.mutate({
															id: todo.id,
															updates: { completed: false },
														})
													}
												>
													<Undo2 className="h-4 w-4" />
												</Button>
											) : (
												<Button
													size="sm"
													onClick={() =>
														updateTodoMutation.mutate({
															id: todo.id,
															updates: { completed: true },
														})
													}
												>
													<CheckCheck className="h-4 w-4" />
												</Button>
											)}
										</>
									)}

									{/* EDIT MODE BUTTONS */}
									{isEditing && (
										<>
											<Button
												variant="default"
												size="sm"
												onClick={() => {
													updateTodoMutation.mutate({
														id: todo.id,
														updates: { title: editValue.trim() },
													});
													setEditingId(null);
												}}
												disabled={!isChanged}
											>
												<Check className="h-4 w-4" />
											</Button>

											<Button
												variant="outline"
												size="sm"
												onClick={() => setEditingId(null)}
											>
												<X className="h-4 w-4" />
											</Button>
										</>
									)}
								</div>
							</motion.div>
						);
					})}
				</AnimatePresence>
			</div>
		</div>
	);
};
