CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION mafiyrm.generate_id()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.id IS NULL THEN
    NEW.id = uuid_generate_v5(
      uuid_ns_dns(),
      CONCAT(
        to_char(NEW.create_date, 'YYYY:MM:DD:HH24:MI:SS:MS:MU'),
        encode(digest(random()::text, 'sha256'), 'escape')
      )
    );
  END IF;

  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION mafiyrm.update_last_update_date_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.last_update_date = now();
  RETURN NEW;
END;
$$ language 'plpgsql';
