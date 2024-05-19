CREATE PROCEDURE `find_paginated_posts`(
    IN `logged_user_id` VARCHAR(191),
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
        post.text AS text,
        post.location AS post_location,
        post.category AS post_category,
        post.reward AS post_reward,
        post.lost_found AS post_lost_found,
        post.animal_type AS post_animal_type,
        post.animal_size AS post_animal_size,
        post.shares_count AS shares,
        post.user_id AS post_author_id,
        post.created_at AS created_at,
        doc.path AS post_media,
        doc.type AS post_media_type,
        doc.mime_type AS post_mime_type,
        usr.id AS author_id,
        usr.name AS post_author,
        usr.user_name AS post_author_username,
        usr_doc.path AS post_author_avatar,
        usr_doc.type AS post_author_avatar_type,
        usr_doc.mime_type AS post_author_avatar_mime_type,
        (SELECT COUNT(*) FROM comment WHERE comment.post_id = post.id) AS comments, 
        (SELECT COUNT(*) FROM interaction_likes WHERE interaction_likes.like_type = "post" AND interaction_likes.post_id = post.id) AS likes,
        CASE
			WHEN usr.id = ''', logged_user_id, ''' THEN true
			    ELSE false
		END AS is_own_post,
        COUNT(*) OVER() AS total_records
    FROM post
        INNER JOIN user usr ON post.user_id = usr.id
        INNER JOIN document doc ON post.id = doc.post_id
        LEFT JOIN document usr_doc ON usr.id = usr_doc.owner_id AND usr_doc.type = "profile_picture"
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

    SET @query = CONCAT(@query, ' ORDER BY post.created_at DESC LIMIT ', result_limit, ' OFFSET ', result_offset);

    PREPARE finalQuery FROM @query;
    EXECUTE finalQuery;
    DEALLOCATE PREPARE finalQuery;
END;