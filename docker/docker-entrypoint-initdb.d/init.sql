CREATE SCHEMA IF NOT EXISTS sch_products;
CREATE USER usr_product WITH PASSWORD '12345';
GRANT USAGE ON SCHEMA sch_products TO usr_product;
