CREATE TABLE IF NOT EXISTS snippets (
    id          BIGSERIAL PRIMARY KEY,
    snippet_id  VARCHAR(255) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    code        TEXT NOT NULL,
    language    VARCHAR(50) NOT NULL,
    filename    VARCHAR(255) NOT NULL,
    description TEXT,
    category    VARCHAR(100),
    is_public   BOOLEAN NOT NULL DEFAULT TRUE,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    views       BIGINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_snippets_snippet_id ON snippets(snippet_id);
CREATE INDEX IF NOT EXISTS idx_snippets_language ON snippets(language);
CREATE INDEX IF NOT EXISTS idx_snippets_category ON snippets(category);
CREATE INDEX IF NOT EXISTS idx_created_at ON snippets(created_at DESC);
