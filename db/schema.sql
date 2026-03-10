CREATE TABLE mcdonalds (
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    image_url     TEXT,
    calories      TEXT,
    calories_unit TEXT,
    protein       TEXT,
    protein_unit  TEXT,
    carbs         TEXT,
    carbs_unit    TEXT,
    fat           TEXT,
    fat_unit      TEXT,
    category      TEXT
);


CREATE TABLE subway (
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    image_url     TEXT,
    calories      TEXT,
    calories_unit TEXT,
    protein       TEXT,
    protein_unit  TEXT,
    carbs         TEXT,
    carbs_unit    TEXT,
    fat           TEXT,
    fat_unit      TEXT,
    category      TEXT
);
