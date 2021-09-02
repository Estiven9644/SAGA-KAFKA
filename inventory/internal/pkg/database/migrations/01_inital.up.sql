CREATE TABLE IF NOT EXISTS inventory (
  product_id int primary key,
  quantity int
);

CREATE TABLE IF NOT EXISTS transactions (
  order_id int,
  product_id int,

  PRIMARY KEY (order_id, product_id)
);
