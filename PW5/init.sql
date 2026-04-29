USE powerdb;

CREATE TABLE systems (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    source VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL
);

CREATE TABLE meters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    system_id INT NOT NULL,
    value INT NOT NULL,
    FOREIGN KEY (system_id) REFERENCES systems(id)
);

CREATE TABLE incidents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    system_id INT NOT NULL,
    description VARCHAR(255),
    FOREIGN KEY (system_id) REFERENCES systems(id)
);
