-- 000001_create_products_table.up.sql

-- Создаем таблицу products с оптимизированной структурой
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cost DECIMAL(10, 2) NOT NULL CHECK (cost >= 0)
);

-- Индекс для поиска по диапазону цен
CREATE INDEX IF NOT EXISTS products_cost_idx ON products (cost);

-- GIN индекс для полнотекстового поиска по описанию
CREATE INDEX IF NOT EXISTS products_description_gin_idx ON products
USING gin (to_tsvector('english', description));

