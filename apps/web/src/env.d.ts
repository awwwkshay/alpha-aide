// Server-side environment variables
declare global {
	namespace NodeJS {
		interface ProcessEnv {
			readonly NITRO_PORT: string;
			readonly NITRO_API_BASE_URL: string;
		}
	}
}

export {};
