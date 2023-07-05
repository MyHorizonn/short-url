CREATE TABLE IF NOT EXISTS urls (
    id SERIAL,
    original varchar(256) NOT NULL,
    short varchar(256) NOT NULL,
    urlkey BIGINT NOT NULL,
    PRIMARY KEY (id)

);