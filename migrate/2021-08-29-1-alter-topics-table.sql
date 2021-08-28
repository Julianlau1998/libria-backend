ALTER TABLE topics
    ADD COLUMN reported BOOLEAN;

ALTER TABLE answers
    ADD COLUMN reported BOOLEAN;

ALTER TABLE comments
    ADD COLUMN reported BOOLEAN;