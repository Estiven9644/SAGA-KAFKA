CREATE TABLE IF NOT EXISTS clients (
  id int primary key
);

CREATE TABLE IF NOT EXISTS cities (
  code int primary key,
  delivery_cost int NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id serial primary key,
  client_id int NOT NULL,
  city_code int NOT NULL,
  delivery_address varchar(128) NOT NULL,
  delivery_cost int NOT NULL,

  FOREIGN KEY (client_id) REFERENCES clients(id),
  FOREIGN KEY (city_code) REFERENCES cities(code)
);

CREATE TABLE IF NOT EXISTS orders_products (
  order_id int NOT NULL,
  product_id int NOT NULL,

  FOREIGN KEY (order_id) REFERENCES orders(id)
)
