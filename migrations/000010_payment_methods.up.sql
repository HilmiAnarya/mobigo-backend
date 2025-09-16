CREATE TABLE payment_methods (
     id SERIAL PRIMARY KEY,
     user_id INT NOT NULL REFERENCES users(id),
     midtrans_token VARCHAR(255) NOT NULL,
     card_type VARCHAR(50),
     masked_card VARCHAR(50),
     is_default BOOLEAN NOT NULL DEFAULT FALSE,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
     deleted_at TIMESTAMP NULL
);
