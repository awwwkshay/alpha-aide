import z from "zod/v4";

export const apiResponseSchema = <T extends z.ZodType>(data: T) =>
	z.object({
		data: data.optional(),
		message: z.string(),
	});

export type IApiResponse<T extends z.ZodType> = z.infer<
	ReturnType<typeof apiResponseSchema<T>>
>;
