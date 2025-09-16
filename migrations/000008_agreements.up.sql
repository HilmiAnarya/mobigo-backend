CREATE TYPE payment_type AS ENUM ('full_payment', 'installment');

CREATE TABLE agreements (
    id SERIAL PRIMARY KEY,
    booking_id INT NOT NULL REFERENCES bookings(id),
    agreement_date TIMESTAMP NOT NULL,
    final_price DECIMAL(10,2) NOT NULL,
    payment_type payment_type NOT NULL,
    terms TEXT,
    signed_by_user BOOLEAN NOT NULL DEFAULT FALSE,
    signed_by_staff BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
