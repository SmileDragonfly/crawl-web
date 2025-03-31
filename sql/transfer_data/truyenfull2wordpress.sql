CREATE TABLE wordpress_db.story_infos LIKE truyenfull.story_infos;
INSERT INTO wordpress_db.story_infos
SELECT * FROM truyenfull.story_infos;

CREATE TABLE wordpress_db.chapters LIKE truyenfull.chapters;
INSERT INTO wordpress_db.chapters
SELECT * FROM truyenfull.chapters;