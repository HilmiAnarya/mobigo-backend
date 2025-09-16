CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  agreement_id INT NOT NULL REFERENCES agreements(id),
  amount DECIMAL(10,2) NOT NULL,
  payment_method VARCHAR(50) NOT NULL,
  status VARCHAR(50) NOT NULL,
  midtrans_transaction_id VARCHAR(100),
  payment_url VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL
);
