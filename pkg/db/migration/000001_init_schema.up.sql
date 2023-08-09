CREATE TABLE "addresss" (
    id bigserial PRIMARY KEY,
    zip_code VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    state VARCHAR NOT NULL,
    country VARCHAR NOT NULL,
    street VARCHAR NOT NULL,
    house_number VARCHAR NOT NULL
);

CREATE TABLE "users" (
    id bigserial PRIMARY KEY,
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    user_type VARCHAR NOT NULL,
    address_id BIGINT NOT NULL,
    orders TEXT[],
    user_cart TEXT[],
    FOREIGN KEY (address_id) REFERENCES "addresss" (id) -- Исправлено здесь
);

CREATE TABLE "products" (
    id bigserial PRIMARY KEY,
    name VARCHAR NOT NULL,
    price DECIMAL NOT NULL,
    description TEXT NOT NULL,
    available_quantity BIGINT NOT NULL,
    category VARCHAR NOT NULL,
    images TEXT[],
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE TABLE "orders" (
    id bigserial PRIMARY KEY,
    order_cart TEXT[],
    total_price DECIMAL NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE "products_to_orders" (
    order_id BIGINT REFERENCES "orders" (id), -- Исправлено здесь
    product_id BIGINT REFERENCES "products" (id), -- Исправлено здесь
    name VARCHAR NOT NULL,
    price DECIMAL NOT NULL,
    buy_quantity BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP,
    PRIMARY KEY (order_id, product_id)
);
