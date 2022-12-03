CREATE TABLE IF NOT EXISTS youtuber (
	`id` TEXT PRIMARY KEY,
   	`title` TEXT NOT NULL,
	`description` TEXT NOT NULL,
	`custom_url` TEXT NOT NULL,
	`default_thumb` TEXT NOT NULL,
	`medium_thumb` TEXT NOT NULL,
	`high_thumb` TEXT NOT NULL,
	`country` TEXT NOT NULL,
	`token` TEXT NOT NULL,
	`refresh_token` TEXT NOT NULL,	
	`publish_at` INTEGER NOT NULL,
	`create_at` INTEGER NOT NULL,
	`update_at` INTEGER NOT NULL
);

INSERT INTO youtuber (
	`id`,
	`title`,
	`description`,
	`custom_url`,
	`default_thumb`,
	`medium_thumb`,
	`high_thumb`,
	`country`,
	`token`,
	`refresh_token`,
	`publish_at`,
	`create_at`,
	`update_at`
) VALUES (
	:id,
	:title,
	:description,
	:custom_url,
	:default_thumb,
	:medium_thumb,
	:high_thumb,
	:country,
	:token,
	:refresh_token,
	:publish_at,
	:create_at,
	:update_at
)