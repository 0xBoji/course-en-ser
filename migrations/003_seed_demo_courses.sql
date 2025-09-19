-- Seed demo courses
INSERT INTO courses (title, description, difficulty) VALUES
(
    'Introduction to Go Programming',
    'Learn the fundamentals of Go programming language, including syntax, data types, functions, and basic concurrency patterns.',
    'Beginner'
),
(
    'Advanced Web Development with React',
    'Master advanced React concepts including hooks, context, state management, and building scalable web applications.',
    'Advanced'
),
(
    'Database Design and SQL',
    'Comprehensive course covering relational database design principles, SQL queries, indexing, and performance optimization.',
    'Intermediate'
),
(
    'Machine Learning Fundamentals',
    'Introduction to machine learning algorithms, data preprocessing, model training, and evaluation techniques using Python.',
    'Intermediate'
)
ON CONFLICT DO NOTHING;
