CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(150) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    is_deleted  BOOLEAN DEFAULT FALSE,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE books (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(200) NOT NULL,
    author      VARCHAR(150) NOT NULL,
    genre       VARCHAR(100),
    price       NUMERIC(10, 2) NOT NULL,
    stock       INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE cart_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    book_id     UUID NOT NULL REFERENCES books(id),
    quantity    INT NOT NULL DEFAULT 1,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, book_id)
);

CREATE TABLE orders (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    total_amount NUMERIC(10, 2) NOT NULL,
    status       VARCHAR(50) DEFAULT 'pending',
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID NOT NULL REFERENCES orders(id),
    book_id     UUID NOT NULL REFERENCES books(id),
    quantity    INT NOT NULL,
    unit_price  NUMERIC(10, 2) NOT NULL
);
