CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        login TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        avatar TEXT,
        date_of_birth TEXT NOT NULL CHECK (
            date_of_birth IS strftime ('%Y-%m-%d', date_of_birth)
            AND date_of_birth >= '1900-01-01'
        ),
        boss_id INTEGER REFERENCES users (id) ON DELETE SET NULL,
        created_at TEXT NOT NULL DEFAULT (strftime ('%Y-%m-%dT%H:%M:%fZ', 'now')),
        updated_at TEXT NOT NULL DEFAULT (strftime ('%Y-%m-%dT%H:%M:%fZ', 'now')),
        token TEXT NOT NULL UNIQUE DEFAULT (lower(hex (randomblob (32)))),
        CHECK (
            boss_id IS NULL
            OR boss_id <> id
        )
    );

CREATE INDEX IF NOT EXISTS users_boss_id ON users (boss_id);

-- Предотвращение изменеие столбца created_at
CREATE TRIGGER IF NOT EXISTS users_protect_created_at BEFORE
UPDATE OF created_at ON users FOR EACH ROW WHEN NEW.created_at IS NOT OLD.created_at BEGIN
SELECT
    RAISE (ABORT, 'Поля created_at нельзя изменить');

END;

-- Предотвращение изменеие столбца created_at
CREATE TRIGGER IF NOT EXISTS users_protect_updated_at BEFORE
UPDATE OF updated_at ON users FOR EACH ROW WHEN NEW.updated_at IS NOT OLD.updated_at
AND NEW.updated_at IS NOT strftime ('%Y-%m-%dT%H:%M:%fZ', 'now') BEGIN
SELECT
    RAISE (ABORT, 'Поле updated_at нельзя изменить вручную');

END;

-- Автоматическое обновление значение updated_at при изменение данных в ряду
CREATE TRIGGER IF NOT EXISTS users_set_updated_at AFTER
UPDATE ON users FOR EACH ROW WHEN NEW.updated_at IS NOT strftime ('%Y-%m-%dT%H:%M:%fZ', 'now') BEGIN
UPDATE users
SET
    updated_at = strftime ('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE
    id = NEW.id;

END;

-- Рекурсивная проверка при создании пользователей и его boss_id
CREATE TRIGGER IF NOT EXISTS users_boss_no_cycle_insert AFTER INSERT ON users FOR EACH ROW WHEN NEW.boss_id IS NOT NULL BEGIN
SELECT
    RAISE (ABORT, 'Циклическое подчинение запрещено')
WHERE
    EXISTS (
        WITH RECURSIVE
            chain (id) AS (
                SELECT
                    NEW.boss_id
                UNION
                SELECT
                    u.boss_id
                FROM
                    users u
                    JOIN chain c ON u.id = c.id
                WHERE
                    u.boss_id IS NOT NULL
            )
        SELECT
            1
        FROM
            chain
        WHERE
            id = NEW.id
    );

END;

-- Рекурсивная проверка при изменении пользователя в его boss_id
CREATE TRIGGER IF NOT EXISTS users_boss_no_cycle_update AFTER
UPDATE OF boss_id ON users FOR EACH ROW WHEN NEW.boss_id IS NOT NULL BEGIN
SELECT
    RAISE (ABORT, 'Циклическое подчинение запрещено')
WHERE
    EXISTS (
        WITH RECURSIVE
            chain (ID) AS (
                SELECT
                    NEW.boss_id
                UNION
                SELECT
                    u.boss_id
                FROM
                    users u
                    JOIN chain c ON u.id = c.id
                WHERE
                    u.boss_id IS NOT NULL
            )
        SELECT
            1
        FROM
            chain
        WHERE
            id = NEW.id
    );

END;

-- Автообновление токена при смене столбцов login или password
CREATE TRIGGER IF NOT EXISTS users_rotate_token AFTER
UPDATE OF login,
password ON USERS FOR EACH ROW WHEN NEW.login IS NOT OLD.login
OR NEW.password IS NOT OLD.password BEGIN
UPDATE users
SET
    token = lower(hex (randomblob (32)))
WHERE
    id = NEW.id;

END;