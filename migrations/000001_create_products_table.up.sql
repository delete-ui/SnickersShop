-- 000001_create_products_table.up.sql

-- Создаем таблицу products с оптимизированной структурой
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cost DECIMAL(10, 2) NOT NULL CHECK (cost >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индекс для поиска по диапазону цен
CREATE INDEX IF NOT EXISTS products_cost_idx ON products (cost);

-- GIN индекс для полнотекстового поиска по описанию
CREATE INDEX IF NOT EXISTS products_description_gin_idx ON products
USING gin (to_tsvector('english', description));

-- Функция для обновления метки времени
CREATE OR REPLACE FUNCTION update_products_updated_at()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматического обновления updated_at при изменении записи
CREATE OR REPLACE TRIGGER products_updated_at_trigger
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_products_updated_at();