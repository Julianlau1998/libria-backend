CREATE TABLE answers (
    answer_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    UserID VARCHAR(30),
    Username VARCHAR(255),
    answer varchar(255) NOT NULL,
    created_date TIMESTAMP NULL,
    updated_date TIMESTAMP NULL,
    PRIMARY KEY(answer_id),
    FOREIGN KEY(topic_id) 
        REFERENCES topics(topic_id)
            ON DELETE CASCADE
);