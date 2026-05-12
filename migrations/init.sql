CREATE TABLE IF NOT EXISTS notifications (
    id           SERIAL PRIMARY KEY,
    user_id      BIGINT       NOT NULL,
    message      TEXT         NOT NULL,
    type         VARCHAR(20)  NOT NULL,
    status       VARCHAR(20)  NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    retry_count  INT          NOT NULL DEFAULT 0
);