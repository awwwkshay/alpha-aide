import { serve } from "@hono/node-server";

import { ENV_VARS } from "@/config";
import { auth, createHonoApp } from "@/lib";

function main() {
	// Initialize the Hono app
	const app = createHonoApp();
	app.on(["POST", "GET"], "/*", (c) => auth.handler(c.req.raw));

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
