CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    balance DECIMAL(20,2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_wallet UNIQUE (user_id)
);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON wallets
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();