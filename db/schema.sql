-- Note that Vitess doesn't support foreign keys, so we simply index column

DROP TABLE IF EXISTS user;
CREATE TABLE user (
    email TEXT PRIMARY KEY NOT NULL,
);

DROP TABLE IF EXISTS locker;
CREATE TABLE locker (
    id INTEGER PRIMARY KEY NOT NULL,
);

DROP TABLE IF EXISTS registration;
CREATE TABLE registration (
    locker INTEGER PRIMARY KEY NOT NULL,
    user TEXT INDEX NOT NULL, -- references user (email) 
    name TEXT NOT NULL,
    expiry TEXT NOT NULL,
    expiryEmailSent TEXT DEFAULT NULL,
);
