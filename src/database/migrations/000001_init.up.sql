CREATE TABLE radius_user_types (
    id BIGSERIAL PRIMARY KEY,
    type_name VARCHAR(64) UNIQUE NOT NULL,  -- e.g., 'student', 'staff', 'guest'
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE radius_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password_hash VARCHAR(128) NOT NULL,   
    user_type_id BIGINT REFERENCES radius_user_types(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE radius_policies (
    id BIGSERIAL PRIMARY KEY,
    user_type_id BIGINT REFERENCES radius_user_types(id) ON DELETE CASCADE,
    attribute VARCHAR(64) NOT NULL,      -- e.g., 'Framed-IP-Address', 'Session-Timeout'
    op VARCHAR(2) DEFAULT ':=',          -- RADIUS operator (:=, ==, +=, etc.)
    value VARCHAR(255) NOT NULL,         -- attribute value
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE radius_nas (
    id BIGSERIAL PRIMARY KEY,
    nas_name VARCHAR(128),               -- description / hostname
    ip_address INET NOT NULL UNIQUE,     -- NAS-IP-Address
    secret VARCHAR(64) NOT NULL,         -- shared secret
    nas_type VARCHAR(64) DEFAULT 'other',-- router, switch, wifi, etc.
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE radius_accounting (
    id BIGSERIAL PRIMARY KEY,
    nas_id BIGINT REFERENCES radius_nas(id) ON DELETE CASCADE,
    username VARCHAR(64) NOT NULL,
    session_id VARCHAR(128) NOT NULL,   -- Acct-Session-Id
    nas_ip INET NOT NULL,               -- NAS-IP-Address
    framed_ip INET,                     -- Framed-IP-Address (user IP)
    acct_status_type VARCHAR(32) NOT NULL, -- Start / Stop / Interim-Update
    acct_session_time INT,              -- duration in seconds
    acct_input_octets BIGINT DEFAULT 0, -- bytes in
    acct_output_octets BIGINT DEFAULT 0,-- bytes out
    acct_terminate_cause VARCHAR(64),   -- e.g., User-Request, Lost-Carrier
    start_time TIMESTAMP,               -- when session started
    stop_time TIMESTAMP,                -- when session ended
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX idx_radius_users_username ON radius_users(username);
CREATE INDEX idx_radius_accounting_session ON radius_accounting(session_id);
CREATE INDEX idx_radius_accounting_user ON radius_accounting(username);
