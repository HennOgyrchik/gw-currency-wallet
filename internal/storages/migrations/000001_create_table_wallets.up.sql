CREATE TABLE IF NOT EXISTS wallets (
                                id serial4 NOT NULL,
                                user_id text NOT NULL,
                                cash json NULL,
                                CONSTRAINT wallets_id_key UNIQUE (id),
                                CONSTRAINT wallets_user_id_key UNIQUE (user_id)
);