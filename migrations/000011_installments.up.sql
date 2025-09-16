ALTER TABLE `agreements`
    ADD COLUMN `payment_type` ENUM('full_payment', 'installment') NOT NULL AFTER `final_price`;