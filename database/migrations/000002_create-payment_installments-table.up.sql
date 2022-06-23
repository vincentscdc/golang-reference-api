CREATE TYPE "payment_installment_status" AS ENUM (
    'pending',
    'paid',
    'due'
);

CREATE TABLE "payment_installments" (
    "id" uuid PRIMARY KEY,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "currency" currency not null,
    "amount" decimal(32, 16) check(amount > 0) not null,
    "due_at" timestamp not null,
    "status" payment_installment_status not null,
    "payment_plan_id" uuid not null,
    CONSTRAINT fk_payment_plans
        FOREIGN KEY(payment_plan_id)
        REFERENCES payment_plans(id)
);
