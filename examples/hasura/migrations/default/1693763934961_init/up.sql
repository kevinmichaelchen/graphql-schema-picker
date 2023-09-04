CREATE TABLE "public"."employees" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "full_name" text NOT NULL,
  PRIMARY KEY ("id")
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "public"."airlines" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "public"."manufacturers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "public"."airports" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "code" text NOT NULL,
  PRIMARY KEY ("id")
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "public"."aircrafts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "model_number" text NOT NULL,
  "serial_number" text NOT NULL,
  "airline_id" uuid NOT NULL,
  "manufacturer_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("airline_id") REFERENCES "public"."airlines"("id") ON
  UPDATE
    no ACTION ON DELETE no ACTION,
    FOREIGN KEY ("manufacturer_id") REFERENCES "public"."manufacturers"("id") ON
  UPDATE
    no ACTION ON DELETE no ACTION
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "public"."flights" (
  "id" uuid NOT NULL,
  "departing_airport_id" uuid NOT NULL,
  "arriving_airport_id" uuid NOT NULL,
  "airline_id" uuid NOT NULL,
  "aircraft_id" uuid NOT NULL,
  "expected_duration_mins" integer NOT NULL,
  "expected_departure_time" timestamptz NOT NULL,
  "expected_arrival_time" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("airline_id") REFERENCES "public"."airlines"("id") ON
  UPDATE
    no ACTION ON DELETE no ACTION,
    FOREIGN KEY ("aircraft_id") REFERENCES "public"."aircrafts"("id") ON
  UPDATE
    no ACTION ON DELETE no ACTION,
    FOREIGN KEY ("departing_airport_id") REFERENCES "public"."airports"("id") ON
  UPDATE
    no ACTION ON DELETE no ACTION,
    FOREIGN KEY ("arriving_airport_id") REFERENCES "public"."airports"("id") ON
  UPDATE
    no ACTION ON DELETE no ACTION
);

CREATE TABLE "public"."staff_roles" (
  "value" text NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("value")
);

INSERT INTO
  "public"."staff_roles"("value", "description")
VALUES
  (E 'PILOT', E 'Pilot');

INSERT INTO
  "public"."staff_roles"("value", "description")
VALUES
  (E 'CO_PILOT', E 'Co-Pilot');

INSERT INTO
  "public"."staff_roles"("value", "description")
VALUES
  (E 'FLIGHT_ATTENDANT', E 'Flight Attendant');

CREATE TABLE "public"."flight_staff" (
  "flight_id" uuid NOT NULL,
  "employee_id" uuid NOT NULL,
  "role" text NOT NULL,
  PRIMARY KEY ("flight_id", "employee_id"),
  FOREIGN KEY ("role") REFERENCES "public"."staff_roles"("value") ON
  UPDATE
    no ACTION ON DELETE no ACTION
);
