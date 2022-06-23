CREATE TYPE "currency" AS ENUM (
    'usdc'
);

CREATE TYPE "payment_status" AS ENUM (
    'pending',
    'complete'
);

CREATE TABLE "payment_plans" (
    "id" uuid PRIMARY KEY,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "currency" currency not null,
    "user_id" uuid not null,
    "amount" decimal(32, 16) not null check(amount > 0),
    "status" payment_status not null
);
