-- =============================================================
-- Новые таблицы для игровых сессий (запустить в pgAdmin)
-- =============================================================

-- Таблица игровых сессий (комнат).
-- Каждая сессия привязана к конкретной игре (game) и имеет
-- уникальный 6-значный код для подключения игроков.
CREATE TABLE game_sessions (
    id                  SERIAL      PRIMARY KEY,
    game_id             INT         NOT NULL REFERENCES game(id) ON DELETE CASCADE,
    code                VARCHAR(6)  NOT NULL UNIQUE,
    status              VARCHAR(20) NOT NULL DEFAULT 'waiting'
                            CHECK (status IN ('waiting', 'active', 'finished')),
    current_question_id INT         REFERENCES questions(id) ON DELETE SET NULL,
    buzzing_player_id   INT         REFERENCES players(id)   ON DELETE SET NULL,
    question_status     VARCHAR(20) NOT NULL DEFAULT 'idle'
                            CHECK (question_status IN ('idle', 'open', 'buzzing')),
    created_at          TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- Таблица игроков в сессии.
-- Хранит кто подключился к какой сессии и их очки внутри этой сессии.
CREATE TABLE session_players (
    id         SERIAL    PRIMARY KEY,
    session_id INT       NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    player_id  INT       NOT NULL REFERENCES players(id)       ON DELETE CASCADE,
    score      INT       NOT NULL DEFAULT 0,
    joined_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (session_id, player_id)
);
