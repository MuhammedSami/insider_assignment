CREATE TYPE message_status AS ENUM (
    'pending',
    'sent',
    'failed',
    'permanent_fail' -- this might make sense if retry doesn't work several times
);

CREATE TABLE messages (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- why I picked UUID, well it simply looked safer
    content VARCHAR(1000),
    recipient_phone_number VARCHAR(20),
    status message_status DEFAULT 'pending' NOT NULL,
    failed_count INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);