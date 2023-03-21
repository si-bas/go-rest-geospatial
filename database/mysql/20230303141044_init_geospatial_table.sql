-- +goose Up
-- +goose StatementBegin
CREATE TABLE geospatial (
    `id` INT NOT NULL AUTO_INCREMENT,
    `gadm_id` VARCHAR(255) NOT NULL UNIQUE,
    `parent_gadm_id` VARCHAR(255) NULL,
    `name` VARCHAR(255) NOT NULL,
    `type` VARCHAR(255) NOT NULL,
    `level` TINYINT(1) UNSIGNED NOT NULL CHECK (
        `level` BETWEEN 1
        AND 5
    ),
    `geometry` MULTIPOLYGON NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    SPATIAL INDEX (`geometry`),
    KEY `idx_gadm_id` (`gadm_id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_updated_at` (`updated_at`)
) ENGINE = InnoDB;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE geospatial;

-- +goose StatementEnd