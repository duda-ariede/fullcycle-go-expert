CREATE TABLE IF NOT EXISTS orders (
    id varchar(255) DEFAULT NULL,
    price float(10,2) DEFAULT NULL,
    tax float(10,2) DEFAULT NULL,
    final_price float(10,2) DEFAULT NULL
)