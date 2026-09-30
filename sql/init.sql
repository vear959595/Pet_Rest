DROP TABLE IF EXISTS tasks;



CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

INSERT INTO tasks (title, description, completed) VALUES
                                                      ('изучить golang', 'Пройти базовый курс', 'TRUE'),
                                                       ('Написать REST API', 'Посмотреть видео про REST', 'FALSE'),
                                                       ('Зарелизить приложение', 'Развернуть приложение на сервере', 'false');