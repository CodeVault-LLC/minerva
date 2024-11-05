SELECT * FROM content LEFT JOIN content_storage ON content.id = content_storage.content_id WHERE scan_id = 1;
