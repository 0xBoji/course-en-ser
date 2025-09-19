-- Add image_url column to courses table for S3 image links
ALTER TABLE courses ADD COLUMN image_url VARCHAR(500);

-- Add index on image_url for potential filtering
CREATE INDEX IF NOT EXISTS idx_courses_image_url ON courses(image_url);

-- Update existing courses with placeholder image URLs (optional)
-- UPDATE courses SET image_url = 'https://your-s3-bucket.s3.amazonaws.com/placeholder-course.jpg' WHERE image_url IS NULL;
