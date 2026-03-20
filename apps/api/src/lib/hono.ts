import type { HttpBindings } from "@hono/node-server";
import { Hono } from "hono";
import { cors } from "hono/cors";

import { ENV_VARS } from "@/config";
import type { ProcessEnv } from "@/schemas";

type Bindings = ProcessEnv & HttpBindings;

type Variables = {};

export const createHonoRouter = () => {
	const router = new Hono<{ Bindings: Bindings; Variables: Variables }>();
	return router;
};

export const createHonoApp = () => {
	const app = createHonoRouter();

	app.use(
		cors({
			origin: ENV_VARS.CLIENT_URLS,
			allowMethods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
			credentials: true,
		}),
	);

	app.get("/health", (c) => {
		return c.json({ status: "ok" });
	});

	return app;
};
