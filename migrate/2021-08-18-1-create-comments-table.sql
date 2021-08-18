CREATE TABLE comments (
    id uuid NOT NULL,
    answer_id uuid NOT NULL,
    UserID VARCHAR(30),
    Username VARCHAR(255),
    comment_text varchar(255) NOT NULL,
    created_date TIMESTAMP NULL,
    updated_date TIMESTAMP NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(answer_id) 
        REFERENCES answers(answer_id)
            ON DELETE CASCADE
);