-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    role          TEXT        NOT NULL DEFAULT 'customer'
                              CHECK (role IN ('admin', 'seller', 'customer')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name      TEXT NOT NULL,
    slug      TEXT NOT NULL UNIQUE,
    parent_id UUID REFERENCES categories(id)
);

CREATE TABLE products (
    id                 UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title              TEXT        NOT NULL,
    description        TEXT        NOT NULL DEFAULT '',
    category_id        UUID        REFERENCES categories(id),
    brand              TEXT        NOT NULL DEFAULT '',
    condition          TEXT        NOT NULL CHECK (condition IN ('excellent', 'good', 'fair')),
    price_cents        INTEGER     NOT NULL CHECK (price_cents > 0),
    seller_id          UUID        REFERENCES users(id),
    instagram_post_url TEXT        NOT NULL DEFAULT '',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE variants (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID        NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size       TEXT        NOT NULL DEFAULT '',
    color      TEXT        NOT NULL DEFAULT '',
    sku        TEXT        UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE inventory (
    variant_id UUID    PRIMARY KEY REFERENCES variants(id) ON DELETE CASCADE,
    quantity   INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved   INTEGER NOT NULL DEFAULT 0 CHECK (reserved >= 0),
    CHECK (quantity >= reserved)
);

CREATE TABLE orders (
    id                UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id           UUID        REFERENCES users(id),
    status            TEXT        NOT NULL DEFAULT 'pending'
                                  CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
    total_cents       INTEGER     NOT NULL CHECK (total_cents > 0),
    stripe_session_id TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id          UUID    PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id    UUID    NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id  UUID    NOT NULL REFERENCES variants(id),
    quantity    INTEGER NOT NULL CHECK (quantity > 0),
    price_cents INTEGER NOT NULL CHECK (price_cents > 0)
);

CREATE TABLE instagram_links (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID        NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    post_url   TEXT        NOT NULL,
    embed_html TEXT        NOT NULL DEFAULT '',
    cached_at  TIMESTAMPTZ
);

CREATE TABLE product_images (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID        NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    r2_key     TEXT        NOT NULL,
    url        TEXT        NOT NULL,
    position   INTEGER     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_created  ON products(created_at DESC);
CREATE INDEX idx_variants_product  ON variants(product_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_product_images_product ON product_images(product_id, position);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS instagram_links;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS variants;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
