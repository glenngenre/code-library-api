CREATE TABLE IF NOT EXISTS rate_limits (
    id          BIGSERIAL PRIMARY KEY,
    ip          VARCHAR(45) NOT NULL UNIQUE,
    blocked_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    unblock_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    reason      VARCHAR(255) DEFAULT 'rate_limit',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_rate_limits_ip ON rate_limits(ip);
CREATE INDEX idx_rate_limits_unblock_at ON rate_limits(unblock_at);
