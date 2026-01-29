ALTER TABLE navigation DROP CONSTRAINT IF EXISTS navigation_single_reference;
ALTER TABLE navigation DROP COLUMN IF EXISTS post_id;
ALTER TABLE navigation DROP COLUMN IF EXISTS category_id;
ALTER TABLE navigation DROP COLUMN IF EXISTS link_source;
