CREATE TABLE
  address (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    pincode INT NOT NULL,
    area TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    country TEXT NOT NULL,
    phone_number BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    last_used TIMESTAMP NOT NULL DEFAULT NOW ()
  );

ALTER TABLE "user" ADD CONSTRAINT fk_user_address FOREIGN KEY (address_id) REFERENCES address (id);