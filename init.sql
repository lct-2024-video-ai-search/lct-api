CREATE TABLE VideoIndex
(
    link              String,                -- ссылка на видеоролик
    audio_description String,                -- расшифровка аудиодорожки из видео
    video_description String,                -- расшифровка видеодорожки из видео
    idx               UInt64,                -- уникальный идентификатор видео в базе данных
    user_description  String,                -- пользовательское описание видео
    created_at        DateTime DEFAULT now() -- время занесения записи в базу данных
)
    Engine = MergeTree
        PRIMARY KEY (idx)