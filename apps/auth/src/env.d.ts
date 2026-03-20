// Server-side environment variables
declare global {
	namespace NodeJS {
		interface ProcessEnv {
			readonly NODE_ENV: string;
			readonly PORT: string;
			readonly DATABASE_URL: string;
			readonly BETTER_AUTH_SECRET: string;
			readonly BETTER_AUTH_URL: string;
			readonly CLIENT_URLS: string;
		}
	}
}

export {};
