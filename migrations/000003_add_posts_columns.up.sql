-- Add columns field to posts table
ALTER TABLE posts 
ADD COLUMN columns INTEGER NOT NULL DEFAULT 1;
