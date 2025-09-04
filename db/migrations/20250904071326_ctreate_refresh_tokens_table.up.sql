CREATE TABLE refresh_tokens (
                                id SERIAL PRIMARY KEY,
                                userId INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token TEXT NOT NULL UNIQUE,
                                createdAt TIMESTAMP NOT NULL DEFAULT now(),
                                expiresAt TIMESTAMP NOT NULL
);