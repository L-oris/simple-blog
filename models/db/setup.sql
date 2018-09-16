CREATE TABLE IF NOT EXISTS Posts (
    ID INT AUTO_INCREMENT NOT NULL, 
    Title VARCHAR(200) NOT NULL, 
    Content TEXT, 
    CreatedAt DATE NOT NULL, 
    PRIMARY KEY (ID)
);

-- MOCK DATA
INSERT INTO Posts (Title, Content) VALUES ("First Post", "First Description");