CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    price int NOT NULL,
    stock int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (category_id) REFERENCES categories (id),
    FOREIGN KEY (shop_id) REFERENCES shops (id)
);