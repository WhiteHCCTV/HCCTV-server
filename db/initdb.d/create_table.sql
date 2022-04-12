CREATE TABLE client (
    `id`      VARCHAR(30)  NOT NULL,
    `pw`      BLOB         NOT NULL,
    `version` VARCHAR(10)  default("0.0.0"),
    PRIMARY KEY (id)
);
CREATE TABLE model (
    `version`   VARCHAR(10)  NOT NULL,
    `weight`    JSON         NOT NULL,
    `accuracy`  INT          NOT NULL,
    PRIMARY KEY (version)
);
