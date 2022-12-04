-- 頻道主資料表
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
	'view' INTEGER NOT NULL,
	'subscriber' INTEGER NOT NULL,
	'video_count' INTEGER NOT NULL,
	`publish_at` INTEGER NOT NULL,
	`create_at` INTEGER NOT NULL,
	`update_at` INTEGER NOT NULL
);
-- yt影片資料表
CREATE TABLE IF NOT EXISTS video (
	`id` TEXT PRIMARY KEY,
   	`channel_id` TEXT NOT NULL,
	`title` TEXT NOT NULL,
	`description` TEXT NOT NULL,
	`default_thumb` TEXT NOT NULL,
	`medium_thumb` TEXT NOT NULL,
	`high_thumb` TEXT NOT NULL,
	`standard_thumb` TEXT NOT NULL,
	`maxres_thumb` TEXT NOT NULL,
	`tags` BLOB NOT NULL,
	`language` TEXT NOT NULL,
	`duration` TEXT NOT NULL,
	`dimension` TEXT NOT NULL,
	`definition` TEXT NOT NULL,
	`caption` INTEGER NOT NULL,
	`view` INTEGER NOT NULL,
	`like` INTEGER NOT NULL,
	`dislike` INTEGER NOT NULL,
	`comment` INTEGER NOT NULL,
	`publish_at` INTEGER NOT NULL,
	`create_at` INTEGER NOT NULL,
	`update_at` INTEGER NOT NULL
);