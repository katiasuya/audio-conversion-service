   -- to create a user
   
    INSERT INTO converter."user" (username, password) VALUES
    ('username', 'password');

   -- to upload an audio
   
    INSERT INTO converter.audio (user_id, name, format, location) VALUES
    ('61e72557-e5af-4bc2-b19e-b1e4c7820d14','song name', 'MP3', 'location');
    
   --to list all the user's audios
    
    SELECT *
    FROM converter.audio a 
    WHERE a.user_id = '61e72557-e5af-4bc2-b19e-b1e4c7820d14';

   -- to get a particular user's audio
    
    SELECT *
    FROM converter.audio a 
    WHERE a.user_id = '61e72557-e5af-4bc2-b19e-b1e4c7820d14'
    AND a.id = '2a4159de-9f06-4920-a9f6-6f612fd0acf5';

   -- to make a request for an audio
  
    INSERT INTO converter.request (original_id, converted_id, status) VALUES
    ('2a4159de-9f06-4920-a9f6-6f612fd0acf5', NULL,'queued');

   -- to get the request history 
   
    SELECT a.name original, r.created, r.updated, r.status
    FROM converter.request r
    INNER JOIN converter.audio a
    ON a.id = r.original_id
    INNER JOIN converter."user" u
    ON a.user_id = u.id
    WHERE u.id='61e72557-e5af-4bc2-b19e-b1e4c7820d14';

   -- to get the request info of an audio

    SELECT a.name original, r.created, r.updated, r.status
    FROM converter.request r
    INNER JOIN converter.audio a
    ON a.id = r.original_id
    INNER JOIN converter."user" u
    ON a.user_id = u.id
    WHERE u.id='61e72557-e5af-4bc2-b19e-b1e4c7820d14'
    AND a.id ='2a4159de-9f06-4920-a9f6-6f612fd0acf5';

   -- to get the list of users sorted by the number of requests

    SELECT u.username, COUNT(r.id) count
    FROM converter.request r
    INNER JOIN converter.audio a
    ON a.id = r.original_id
    RIGHT JOIN converter."user" u
    ON a.user_id = u.id
    GROUP BY u.username 
    ORDER BY count DESC;

   -- to get the user with the most number of conversion requests

    SELECT u.username
    FROM converter.request r
    INNER JOIN converter.audio a
    ON a.id = r.original_id
    INNER JOIN converter."user" u
    ON a.user_id = u.id
    GROUP BY u.username
    HAVING COUNT(r.id) = 
    (SELECT MAX(s1.c)
    FROM (SELECT COUNT(r.id) c
    FROM converter.request r
    INNER JOIN converter.audio a
    ON a.id = r.original_id
    INNER JOIN converter."user" u
    ON a.user_id = u.id
    GROUP BY u.username) s1);

   -- toget the list of users which did not request to convert from mp3

    SELECT  u.username
	 FROM converter."user" u 
    WHERE u.username NOT IN 
    (SELECT u.username
    FROM  converter.request r
    INNER JOIN converter.audio a
    ON a.id = r.original_id 
    INNER JOIN converter."user" u
    ON u.id = a.user_id 
    WHERE a.format='MP3');


