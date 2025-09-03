CREATE TABLE userProfiles (
    id SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    firstName VARCHAR(100),
    lastName VARCHAR(100),
    middleName VARCHAR(100),
    bio TEXT,
    avatarUrl VARCHAR(255),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_profiles_user_id ON userProfiles(userId);