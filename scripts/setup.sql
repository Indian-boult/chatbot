CREATE TABLE IF NOT EXISTS images (
                                      id INT AUTO_INCREMENT PRIMARY KEY,
                                      identifier VARCHAR(255) NOT NULL,
                                      url VARCHAR(255) NOT NULL,
                                      UNIQUE KEY unique_identifier (identifier)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
