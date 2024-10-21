CREATE TABLE bookings
(
    id             uuid PRIMARY KEY,
    first_name     VARCHAR(255) NOT NULL,
    last_name      VARCHAR(255) NOT NULL,
    gender         VARCHAR(50)  NOT NULL,
    birthday       VARCHAR(255) NOT NULL,

    launch_pad_id  VARCHAR(255) NOT NULL,
    destination_id VARCHAR(255) NOT NULL,
    launch_date    TIMESTAMPTZ  NOT NULL,

    created_at     TIMESTAMPTZ  NOT NULL,
    updated_at     TIMESTAMPTZ  NOT NULL
);