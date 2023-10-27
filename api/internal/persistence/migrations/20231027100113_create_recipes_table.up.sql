CREATE TABLE recipes(
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    ingredients TEXT[],

    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);