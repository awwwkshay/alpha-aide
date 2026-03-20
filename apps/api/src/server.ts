import { serve } from "@hono/node-server";

import { ENV_VARS } from "@/config";
import { todosRouter } from "@/controllers";
import { createHonoApp } from "@/lib";

function main() {
	// Initialize the Hono app
	const app = createHonoApp();
	app.route("/todos", todosRouter);

	// Start the server
	const server = serve({ fetch: app.fetch, port: ENV_VARS.PORT }, (info) => {
		console.log(`Server is running on http://localhost:${info.port}`);
	});

	// Graceful shutdown
	process.on("SIGINT", () => {
		server.close();
		process.exit(0);
	});
	process.on("SIGTERM", () => {
		server.close((err) => {
			if (err) {
				console.error(err);
				process.exit(1);
			}
			process.exit(0);
		});
	});
}

main();
