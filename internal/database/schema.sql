CREATE TABLE IF NOT EXISTS locker (
    id varchar(255) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS registration (
    locker varchar(255) NOT NULL,
    user varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    expiry datetime NOT NULL,
    expiryEmailSent boolean DEFAULT FALSE,
    PRIMARY KEY (locker)
);

CREATE INDEX IF NOT EXISTS user_registration 
ON registration (user);
