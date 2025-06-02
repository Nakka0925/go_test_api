USE test_db;

CREATE TABLE IF NOT EXISTS  html_table (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255),
    content TEXT
);