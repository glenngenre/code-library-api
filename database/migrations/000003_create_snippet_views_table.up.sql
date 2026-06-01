CREATE TABLE IF NOT EXISTS snippet_views (
    snippet_id  VARCHAR(36) NOT NULL,
    hash        VARCHAR(64) NOT NULL,
    CONSTRAINT uq_snippet_views UNIQUE (snippet_id, hash),
    CONSTRAINT fk_snippet_views_snippet FOREIGN KEY (snippet_id) REFERENCES snippets(snippet_id) ON DELETE CASCADE
);

CREATE INDEX idx_snippet_views_snippet_id ON snippet_views(snippet_id);
