CREATE TABLE IF NOT EXISTS Posts (
    ID INT AUTO_INCREMENT NOT NULL, 
    Title VARCHAR(200) NOT NULL, 
    Description TEXT, 
    CreatedAt DATE NOT NULL, 
    PRIMARY KEY (ID)
);

-- MOCK DATA
INSERT INTO Posts (Title, Description) VALUES ("First Post", "First Description");