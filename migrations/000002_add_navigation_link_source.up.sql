-- Add link_source column to navigation table
ALTER TABLE navigation 
ADD COLUMN link_source TEXT CHECK (link_source IN ('custom', 'category', 'post')) DEFAULT 'custom';

-- Add reference columns for category and post links
ALTER TABLE navigation 
ADD COLUMN category_id UUID REFERENCES categories(id) ON DELETE SET NULL;

ALTER TABLE navigation 
ADD COLUMN post_id UUID REFERENCES posts(id) ON DELETE SET NULL;

-- Add constraint: only one reference can be set
ALTER TABLE navigation 
ADD CONSTRAINT navigation_single_reference 
CHECK (
    (link_source = 'custom' AND category_id IS NULL AND post_id IS NULL) OR
    (link_source = 'category' AND category_id IS NOT NULL AND post_id IS NULL) OR
    (link_source = 'post' AND post_id IS NOT NULL AND category_id IS NULL)
);
