import { createFileRoute } from "@tanstack/react-router";

import { TodoPanel } from "@/components";

export const Route = createFileRoute("/dashboard/")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<div className="h-full flex place-content-center">
			<TodoPanel />
		</div>
	);
}
