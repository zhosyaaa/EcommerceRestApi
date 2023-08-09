
ALTER TABLE users DROP CONSTRAINT users_address_id_fkey;
ALTER TABLE products_to_orders DROP CONSTRAINT products_to_orders_order_id_fkey;
ALTER TABLE products_to_orders DROP CONSTRAINT products_to_orders_product_id_fkey;

DROP TABLE IF EXISTS products_to_orders;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS addresss;
