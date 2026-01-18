CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    balance INT DEFAULT 0,
    isAdmin BOOL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    userId BIGINT NOT NULL,
	adminEmail CHAR(255),
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount INT NOT NULL
);

CREATE TABLE IF NOT EXISTS admins (
    email CHAR(255) NOT NULL PRIMARY KEY
);

DELIMITER $$
CREATE TRIGGER IF NOT EXISTS update_user_balance AFTER INSERT ON transactions
    FOR EACH ROW
BEGIN
    -- Update the user's balance
    UPDATE users SET balance = balance + NEW.amount WHERE id = NEW.userId;
END$$
DELIMITER ;
