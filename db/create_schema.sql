CREATE TABLE users (
    id BIGINT NOT NULL PRIMARY KEY,
    username CHAR(255) NOT NULL,
    balance INT DEFAULT 0
);

CREATE TABLE transactions (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    userId BIGINT NOT NULL,
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount INT NOT NULL
);

CREATE TABLE admins (
    email CHAR(255) NOT NULL PRIMARY KEY
);

DELIMITER $$
CREATE TRIGGER update_user_balance AFTER INSERT ON transactions
    FOR EACH ROW
BEGIN
    -- Update the user's balance
    UPDATE users SET balance = balance + NEW.amount WHERE id = NEW.userId;
END$$
DELIMITER ;
