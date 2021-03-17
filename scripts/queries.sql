   -- to create a user
   
    INSERT INTO converter."user" (username, password) VALUES
    ('username', 'password');

   -- to log in the user 

    SELECT password
    FROM converter."user" 
    WHERE username = 'username';

   -- to upload an audio and create a request
   
    WITH audio_id AS (
            INSERT INTO converter.audio (name, format, location) VALUES
            ('song name', 'mp3', 'location') RETURNING id)

    INSERT INTO converter.request (user_id, source_id, source_format, target_id, target_format, status)
    SELECT '12feccec-3974-4dc2-ac63-b4838c7bf0eb', id, 'mp3', NULL, 'wav', 'queued'
    FROM audio_id;

   -- to list all the user's audios
    
    SELECT a.name
    FROM converter.audio a 
    JOIN converter.request r
    ON a.id=r.source_id OR a.id=r.target_id 
    WHERE r.user_id = '61e72557-e5af-4bc2-b19e-b1e4c7820d14';

   -- to get a particular user's audio
    
    SELECT *
    FROM converter.audio a 
    JOIN converter.request r
    ON a.id=r.source_id OR a.id=r.target_id 
    WHERE r.user_id = '61e72557-e5af-4bc2-b19e-b1e4c7820d14'
    AND a.id = '2a4159de-9f06-4920-a9f6-6f612fd0acf5';

   -- to get the request history 
   
    SELECT r.id, a.name, source_format, target_format, r.created, r.updated, r.status
    FROM converter.request r
    JOIN converter.audio a
    ON a.id = r.source_id
    WHERE r.user_id='61e72557-e5af-4bc2-b19e-b1e4c7820d14';

   -- to get the request info of an audio

    SELECT a.name, source_format, target_format, r.created, r.updated, r.status
    FROM converter.request r
    JOIN converter.audio a
    ON a.id = r.source_id
    WHERE r.user_id='61e72557-e5af-4bc2-b19e-b1e4c7820d14'
    AND a.id ='2a4159de-9f06-4920-a9f6-6f612fd0acf5';

   -- to get the list of users sorted by the number of requests

    SELECT u.username, COUNT(r.id) count
    FROM converter.request r
    RIGHT JOIN converter."user" u
    ON r.user_id = u.id
    GROUP BY u.username 
    ORDER BY count DESC;

   -- to get the user with the most number of conversion requests

    SELECT u.username
    FROM converter.request r
    JOIN converter."user" u
    ON r.user_id = u.id
    GROUP BY u.username
    HAVING COUNT(r.id) = (SELECT MAX(count.c)
                          FROM (SELECT COUNT(r.id) c
                                FROM converter.request r
                                JOIN converter."user" u
                                ON r.user_id = u.id
                                GROUP BY u.username) count);

   -- to get the list of users which did not request to convert from mp3

    SELECT  u.username
    FROM converter."user" u 
    WHERE u.username NOT IN (SELECT u.username
                             FROM  converter.request r
                             JOIN converter."user" u
                             ON u.id = r.user_id 
                             JOIN converter.audio a
                             ON a.id = r.source_id 
                             WHERE a.format='mp3');


