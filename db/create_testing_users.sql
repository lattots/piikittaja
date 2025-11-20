CREATE USER 'app_tester'@'%' IDENTIFIED BY 'app_test_pwd';

GRANT SELECT, INSERT, UPDATE, DELETE ON piikki_test_db.users TO 'app_tester'@'%';
GRANT SELECT, INSERT ON piikki_test_db.transactions TO 'app_tester'@'%';

FLUSH PRIVILEGES;
