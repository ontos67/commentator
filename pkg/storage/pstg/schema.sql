DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER DEFAULT 0,
    comment_text TEXT, 
    pub_time INTEGER DEFAULT 0,
    parent_type VARCHAR(1),
    parent_id INTEGER DEFAULT 0
);
