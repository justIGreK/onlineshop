CREATE TABLE "users" (
  "id" serial NOT NULL unique,
  "login" VARCHAR(255) NOT NULL unique,
  "password" VARCHAR(255) NOT NULL,
  "balance" INT NOT NULL DEFAULT '0',
  "is_active" BOOLEAN NOT NULL DEFAULT 'TRUE'
);

CREATE TABLE "cart" (
  "id" serial NOT NULL unique,
  "user_id" int NOT NULL,
  "product_id" int NOT NULL,
  "quantity" int NOT NULL
);

CREATE TABLE "orders" (
  "id" serial NOT NULL unique,
  "user_id" bigint NOT NULL,
  "price_before" float NOT NULL,
  "price_after" float NOT NULL,
  "payed_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "orders_items" (
  "id" serial NOT NULL unique,
  "order_id" int NOT NULL,
  "product_id" int NOT NULL,
  "quantity" int NOT NULL,
  "total_cost" float
);

CREATE TABLE "products" (
  "id" serial NOT NULL unique,
  "name" VARCHAR(255) NOT NULL unique,
  "cost" float NOT NULL,
  "description" VARCHAR(255),
  "amount" int NOT NULL,
  "is_active" BOOLEAN NOT NULL DEFAULT 'TRUE'
);

ALTER TABLE "cart" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "orders_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
