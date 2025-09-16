ALTER TABLE `bookings`
    ADD COLUMN `proposed_datetime` TIMESTAMP NULL DEFAULT NULL AFTER `status`;