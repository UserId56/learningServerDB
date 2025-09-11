CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE refresh_tokens (
                                id SERIAL PRIMARY KEY,
                                userId INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
                                createdAt TIMESTAMP NOT NULL DEFAULT now(),
                                expiresAt TIMESTAMP NOT NULL
);