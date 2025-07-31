DROP TABLE IF EXISTS vector_store;
CREATE TABLE vector_store (
    vector_id VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    embedding VECTOR(1024) NOT NULL,
    CONSTRAINT vector_store_pk PRIMARY KEY (vector_id)
);