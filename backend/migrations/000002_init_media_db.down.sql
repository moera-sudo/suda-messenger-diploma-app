DROP INDEX IF EXISTS idx_media_entity_lookup;
DROP INDEX IF EXISTS idx_media_entity_media;
DROP INDEX IF EXISTS idx_media_bucket_key;
DROP INDEX IF EXISTS idx_media_status;
DROP INDEX IF EXISTS idx_media_owner;

-- 2. Таблица связей (зависит от media)
DROP TABLE IF EXISTS media_entity_links;

-- 3. Основная таблица
DROP TABLE IF EXISTS media_files;
