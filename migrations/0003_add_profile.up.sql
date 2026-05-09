CREATE TABLE IF NOT EXISTS profile (
    id         BIGSERIAL PRIMARY KEY,
    full_name  VARCHAR(255) NOT NULL DEFAULT 'Балашов Максим',
    title      VARCHAR(255) NOT NULL DEFAULT 'Full-Stack разработчик',
    about      TEXT,
    email      VARCHAR(255),
    telegram   VARCHAR(255),
    github     VARCHAR(255),
    phone      VARCHAR(50)
);

INSERT INTO profile (full_name, title, email, telegram, github, phone)
VALUES (
    'Балашов Максим',
    'Full-Stack разработчик',
    '79164211428@yandex.ru',
    'https://t.me/BalashovMaximm',
    'https://github.com/Maxim-Ba',
    '+7 916 421-14-28'
)
ON CONFLICT DO NOTHING;
