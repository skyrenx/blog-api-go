CREATE TABLE blog_entry_sequence (
    next_id INT NOT NULL
);

-- Initialize with starting value 1
INSERT INTO blog_entry_sequence (next_id) VALUES (1);
