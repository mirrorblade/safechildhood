CREATE TABLE IF NOT EXISTS complaints (
    id UUID CONSTRAINT pk PRIMARY KEY,
    coordinates VARCHAR(32) NOT NULL,
    short_description VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);