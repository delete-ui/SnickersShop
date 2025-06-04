-- 000001_create_products_table.down.sql

-- Удаляем таблицу и все зависимые объекты
DROP TABLE IF EXISTS products CASCADE;

-- Удаляем функцию, если она осталась (на случай если CASCADE не сработал)
DROP FUNCTION IF EXISTS update_products_updated_at();