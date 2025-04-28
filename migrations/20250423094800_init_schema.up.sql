-- File: 20250423120000_init_schema.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    is_superuser BOOLEAN
);

CREATE TABLE IF NOT EXISTS messages(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    email TEXT,
    message TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS articles(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT,
    description TEXT, 
    content TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    public BOOLEAN,
    cover_image_id UUID
);

CREATE TABLE IF NOT EXISTS images(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    uploaded_at TIMESTAMP DEFAULT now(),
    article_id UUID,
    CONSTRAINT fk_article FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

ALTER TABLE articles
ADD CONSTRAINT fk_cover_image
FOREIGN KEY (cover_image_id)
REFERENCES images(id)
ON DELETE SET NULL;