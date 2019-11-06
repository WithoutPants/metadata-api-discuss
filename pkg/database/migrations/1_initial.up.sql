CREATE TABLE `tags` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `created_at` datetime not null,
  `updated_at` datetime not null
);
CREATE TABLE `studios` (
  `id` integer not null primary key autoincrement,
  `image` blob,
  `name` varchar(255) not null,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null
);
CREATE TABLE `scenes_tags` (
  `scene_id` integer,
  `tag_id` integer,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`)
);
CREATE TABLE `scenes` (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `details` text,
  `url` varchar(255),
  `date` date,
  `size` varchar(255),
  `duration` float,
  `video_codec` varchar(255),
  `audio_codec` varchar(255),
  `width` tinyint,
  `height` tinyint,
  `framerate` float,
  `bitrate` integer,
  `studio_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);
CREATE TABLE `performers_scenes` (
  `performer_id` integer,
  `scene_id` integer,
  foreign key(`performer_id`) references `performers`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);
CREATE TABLE `performers` (
  `id` integer not null primary key autoincrement,
  `image` blob,
  `name` varchar(255) not null,
  `gender` varchar(255),
  `url` varchar(255),
  `twitter` varchar(255),
  `instagram` varchar(255),
  `birthdate` date,
  `ethnicity` varchar(255),
  `country` varchar(255),
  `eye_color` varchar(255),
  `height` varchar(255),
  `measurements` varchar(255),
  `fake_tits` varchar(255),
  `career_length` varchar(255),
  `tattoos` varchar(255),
  `piercings` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null
);
CREATE TABLE `performer_aliases` (
  `performer_id` integer,
  `alias` varchar(255),
  foreign key(`performer_id`) references `performers`(`id`)
);

CREATE INDEX `index_tags_on_name` on `tags` (`name`);
CREATE INDEX `index_studios_on_name` on `studios` (`name`);
CREATE INDEX `index_scenes_tags_on_tag_id` on `scenes_tags` (`tag_id`);
CREATE INDEX `index_scenes_tags_on_scene_id` on `scenes_tags` (`scene_id`);
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);
CREATE INDEX `index_performers_scenes_on_scene_id` on `performers_scenes` (`scene_id`);
CREATE INDEX `index_performers_scenes_on_performer_id` on `performers_scenes` (`performer_id`);
CREATE INDEX `index_performers_on_name` on `performers` (`name`);
CREATE INDEX `index_performers_on_alias` on `performer_aliases` (`alias`);
