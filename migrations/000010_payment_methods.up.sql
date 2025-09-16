ALTER TABLE `vehicles`
    MODIFY COLUMN `status` ENUM('available', 'booked', 'sold', 'on_installment') NOT NULL DEFAULT 'available';