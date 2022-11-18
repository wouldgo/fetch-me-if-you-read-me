CREATE TABLE IF NOT EXISTS mafiyrm.images (
  id UUID NOT NULL UNIQUE,
  used_in VARCHAR(255) NOT NULL,
  last_update_date TIMESTAMP WITH TIME ZONE,
  create_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (used_in)
);

CREATE INDEX IF NOT EXISTS images_id_idx
  ON mafiyrm.images (id);

CREATE INDEX IF NOT EXISTS images_used_in_idx
  ON mafiyrm.images (used_in);

DROP TRIGGER IF EXISTS update_last_update_date
  ON mafiyrm.images;
CREATE TRIGGER update_last_update_date
  BEFORE UPDATE
  ON mafiyrm.images
  FOR EACH ROW
  EXECUTE PROCEDURE mafiyrm.update_last_update_date_column();

DROP TRIGGER IF EXISTS generate_id ON mafiyrm.images;
CREATE TRIGGER generate_id
  BEFORE INSERT
  ON mafiyrm.images
  FOR EACH ROW
  EXECUTE PROCEDURE mafiyrm.generate_id();

---

CREATE TABLE IF NOT EXISTS mafiyrm.who (
  id UUID NOT NULL UNIQUE,
  remote_addr VARCHAR(255),
  meta JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_update_date TIMESTAMP WITH TIME ZONE,
  create_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (remote_addr)
);

CREATE INDEX IF NOT EXISTS who_id_idx
  ON mafiyrm.who (id);

CREATE INDEX IF NOT EXISTS who_remote_addr_idx
  ON mafiyrm.who (remote_addr);

DROP TRIGGER IF EXISTS update_last_update_date
  ON mafiyrm.who;
CREATE TRIGGER update_last_update_date
  BEFORE UPDATE
  ON mafiyrm.who
  FOR EACH ROW
  EXECUTE PROCEDURE mafiyrm.update_last_update_date_column();

DROP TRIGGER IF EXISTS generate_id ON mafiyrm.who;
CREATE TRIGGER generate_id
  BEFORE INSERT
  ON mafiyrm.who
  FOR EACH ROW
  EXECUTE PROCEDURE mafiyrm.generate_id();

---

CREATE TABLE IF NOT EXISTS mafiyrm.images_accesed (
  image_fk UUID NOT NULL,
  who_fk UUID NOT NULL,
  create_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS images_accesed_image_fk_idx
  ON mafiyrm.images_accesed (image_fk);

CREATE INDEX IF NOT EXISTS images_accesed_who_idx
  ON mafiyrm.images_accesed (who_fk);
