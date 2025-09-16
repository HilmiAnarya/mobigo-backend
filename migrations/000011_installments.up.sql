CREATE TABLE installments (
  id SERIAL PRIMARY KEY,
  payment_id INT NOT NULL REFERENCES payments(id),
  due_date DATE NOT NULL,
  amount_due DECIMAL(10,2) NOT NULL,
  penalty_amount DECIMAL(10,2) DEFAULT 0.00,
  total_due DECIMAL(10,2) NOT NULL,
  status VARCHAR(50) NOT NULL,
  paid_date DATE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL
);
