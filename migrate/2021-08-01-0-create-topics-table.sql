CREATE TABLE topics (
    topic_id uuid NOT NULL,
    UserID VARCHAR(30),
    Username VARCHAR(255),
    title varchar(255) NOT NULL,
    body varchar(255) NULL,
    created_date TIMESTAMP NULL,
    updated_date TIMESTAMP NULL,
    PRIMARY KEY(topic_id)
);