CREATE PROCEDURE `find_paginated_posts`(
    IN `lost_found` VARCHAR(30),
    IN `neighborhood` TEXT,
    IN `reward` CHAR(1),
    IN `user_id` VARCHAR(191),
    IN `result_limit` INT,
    IN `result_offset` INT
)
BEGIN
    SET @query = CONCAT('
    SELECT DISTINCT
        post.id AS post_id,
        post.location AS post_location,
        post.category AS post_category,
        post.reward AS post_reward,
        post.lost_found AS post_lost_found,
        post.shares_count AS shares,
        post.user_id AS post_author_id,
        post.created_at AS created_at,
        doc.path AS post_media,
        doc.type AS post_media_type,
        doc.mime_type AS post_mime_type,
        usr.id AS author_id,
        usr.name AS post_author,
        usr_doc.path AS post_author_avatar,
        usr_doc.type AS post_author_avatar_type,
        usr_doc.mime_type AS post_author_avatar_mime_type,
        (SELECT COUNT(*) FROM comment WHERE comment.post_id = post.id) AS comments, 
        (SELECT COUNT(*) FROM interaction_likes WHERE interaction_likes.like_type = "post" AND interaction_likes.post_id = post.id) AS likes,
        COUNT(*) OVER() AS total_records
    FROM post
        INNER JOIN user usr ON post.user_id = usr.id
        INNER JOIN document doc ON post.id = doc.post_id
        INNER JOIN document usr_doc ON usr.id = doc.owner_id
    WHERE post.deleted_at IS NULL');

    IF lost_found IS NOT NULL THEN
        SET @query = CONCAT(@query, ' AND post.lost_found = ''', lost_found, '''');
    END IF;

    IF neighborhood IS NOT NULL THEN
        SET @query = CONCAT(@query, ' AND (post.location LIKE ''%', neighborhood, '%'')');
    END IF;

    IF reward IS NOT NULL THEN
        SET @query = CONCAT(@query, ' AND post.reward = ''', reward, '''');
    END IF;

    IF user_id IS NOT NULL THEN
        SET @query = concat(@query, ' AND post.user_id = ''', user_id,''' ');
    END IF;

    SET @query = CONCAT(@query, ' LIMIT ', result_limit, ' OFFSET ', result_offset);

    PREPARE finalQuery FROM @query;
    EXECUTE finalQuery;
    DEALLOCATE PREPARE finalQuery;
END;