CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    vehicle_id INT NOT NULL REFERENCES vehicles(id),
    status VARCHAR(50) NOT NULL,
    proposed_datetime TIMESTAMP NOT NULL,
    decline_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
