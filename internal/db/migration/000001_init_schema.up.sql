CREATE TABLE "token" (
    "id" bigserial PRIMARY KEY,
    "access_token" varchar NOT NULL,
    "created_at" timestamptz DEFAULT (now())
);

-- CREATE FUNCTION trigger_set_timestamp()
--     RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.created_at = NOW();
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER set_timestamp
--     BEFORE
--         UPDATE ON token
--     FOR EACH ROW
-- EXECUTE PROCEDURE trigger_set_timestamp();