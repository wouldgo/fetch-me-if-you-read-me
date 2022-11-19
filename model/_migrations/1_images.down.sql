DROP INDEX IF EXISTS mafiyrm.images_id_idx;
DROP INDEX IF EXISTS mafiyrm.images_used_in_idx;
DROP INDEX IF EXISTS mafiyrm.who_id_idx;
DROP INDEX IF EXISTS mafiyrm.who_ip_idx;
DROP INDEX IF EXISTS mafiyrm.images_accessed_image_fk_idx;
DROP INDEX IF EXISTS mafiyrm.images_accessed_who_idx;

DROP TABLE IF EXISTS mafiyrm.images_accessed CASCADE;
DROP TABLE IF EXISTS mafiyrm.who CASCADE;
DROP TABLE IF EXISTS mafiyrm.images CASCADE;
