CREATE TABLE users (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный идентификатор пользователя, UUID хорош для распределённых систем.
       email VARCHAR(255) UNIQUE NOT NULL, -- Основной идентификатор для входа и связи с пользователем. Ограничение уникально
       username VARCHAR(50) UNIQUE NOT NULL, -- Уникальное публичное имя пользователя, может использоваться для отображения или
       password_hash TEXT NOT NULL, -- Хранит хэш пароля пользователя. Никогда не храним пароли в открытом виде. !!! Д
       first_name VARCHAR(100), -- Имя пользователя, для отображения или персонализации.
       last_name VARCHAR(100), -- Фамилия пользователя.
       is_active BOOLEAN DEFAULT TRUE, -- Флаг активности пользователя. Можно деактивировать вместо полного удаления.
       role VARCHAR(50) DEFAULT 'user', -- Роль пользователя, например: 'user', 'admin'. Удобно для простых RBAC (role-bas
       last_login_at TIMESTAMP WITH TIME ZONE, -- Дата и время последнего входа. Полезно для безопасности и аналитики.
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Дата и время регистрации пользователя.
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Дата последнего изменения данных пользователя.
);
