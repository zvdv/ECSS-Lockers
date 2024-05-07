-- Note that Vitess doesn't support foreign keys, so we simply index column

DROP TABLE IF EXISTS user;
CREATE TABLE user (
    email varchar(255) NOT NULL,
    PRIMARY KEY (email)
);

DROP TABLE IF EXISTS locker;
CREATE TABLE locker (
    id varchar(255) NOT NULL,
    PRIMARY KEY(id)
);

DROP PROCEDURE IF EXISTS SeedLockers;
DELIMITER //

CREATE PROCEDURE SeedLockers()
BEGIN
    DECLARE first INT DEFAULT 1;
    DECLARE last INT DEFAULT 200;
    
    WHILE first <= last DO
        INSERT INTO locker VALUES(CONCAT('ELW ', LPAD(first, 3, " ")));
        SET first = first + 1;
    END WHILE;
END//

DELIMITER ;

CALL SeedLockers();


DROP TABLE IF EXISTS registration;
CREATE TABLE registration (
    locker varchar(255) NOT NULL,
    user varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    expiry datetime NOT NULL,
    expiryEmailSent datetime DEFAULT NULL,
    PRIMARY KEY (locker),
    INDEX (user) -- references user (email) 
);
