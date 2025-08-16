CREATE TABLE IF NOT EXISTS orders (
  "order_uid" varchar primary key,
  "track_number" varchar,
  "entry" varchar,
  "locale" varchar,
  "internal_signature" varchar,
  "customer_id" varchar,
  "delivery_service" varchar,
  "shardkey" varchar,
  "sm_id" varchar,
  "date_created" timestamp,
  "oof_shard" varchar
);

CREATE TABLE IF NOT EXISTS delivery (
  "order_uid" varchar primary key,
  "name" varchar,
  "phone" varchar,
  "zip" varchar,
  "city" varchar,
  "address" varchar,
  "region" varchar,
  "email" varchar
);

CREATE TABLE IF NOT EXISTS payment (
  "order_uid" varchar primary key,
  "transaction" varchar,
  "request_id" varchar,
  "currency" varchar,
  "provider" varchar,
  "amount" int,
  "payment_dt" int,
  "bank" varchar,
  "delivery_cost" int,
  "goods_total" int,
  "custom_fee" int
);

CREATE TABLE IF NOT EXISTS items (
  "order_uid" varchar primary key,
  "chrt_id" int,
  "track_number" varchar,
  "price" int,
  "rid" varchar,
  "name" varchar,
  "sale" int,
  "size" varchar,
  "total_price" int,
  "nm_id" int,
  "brand" varchar,
  "status" int
);

ALTER TABLE delivery ADD FOREIGN KEY ("order_uid") REFERENCES orders ("order_uid");

ALTER TABLE payment ADD FOREIGN KEY ("order_uid") REFERENCES orders ("order_uid");

ALTER TABLE items ADD FOREIGN KEY ("order_uid") REFERENCES orders ("order_uid");