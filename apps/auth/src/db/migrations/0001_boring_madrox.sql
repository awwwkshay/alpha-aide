ALTER TABLE "users" DROP CONSTRAINT "users_phone_number_unique";--> statement-breakpoint
ALTER TABLE "users" DROP COLUMN "phone_number";--> statement-breakpoint
ALTER TABLE "users" DROP COLUMN "phone_number_verified";