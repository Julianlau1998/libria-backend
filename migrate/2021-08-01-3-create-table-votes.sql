CREATE TABLE votes (
    vote_id uuid NOT NULL,
    answer_id uuid,
    userID VARCHAR(255),
    upvote VARCHAR(30),
    PRIMARY KEY(vote_id),
    FOREIGN KEY (answer_id)
     REFERENCES answers(answer_id)
     ON DELETE CASCADE
);