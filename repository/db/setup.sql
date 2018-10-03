CREATE TABLE IF NOT EXISTS Posts (
    ID INT AUTO_INCREMENT NOT NULL, 
    Title VARCHAR(200) NOT NULL, 
    Content TEXT, 
    Picture VARCHAR(100) DEFAULT '', 
    CreatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    PRIMARY KEY (ID)
);

-- MOCK DATA
INSERT INTO Posts (Title, Content) VALUES ("First Post", "First Description");