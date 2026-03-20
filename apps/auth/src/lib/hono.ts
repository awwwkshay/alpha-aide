import type { HttpBindings } from "@hono/node-server";
import { Hono } from "hono";
import { cors } from "hono/cors";

import { ENV_VARS } from "@/config";
import { auth } from "@/lib";
import type { ProcessEnv } from "@/schemas";

type Bindings = ProcessEnv & HttpBindings;

type Variables = {
	user: typeof auth.$Infer.Session.user | null;
	session: typeof auth.$Infer.Session.session | null;
};

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

	app.use("*", async (c, next) => {
		const session = await auth.api.getSession({ headers: c.req.raw.headers });
		if (!session) {
			c.set("user", null);
			c.set("session", null);
			await next();
			return;
		}
		c.set("user", session.user);
		c.set("session", session.session);
		await next();
	});

	return app;
};
