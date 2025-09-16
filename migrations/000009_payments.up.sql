ALTER TABLE `bookings`
    ADD COLUMN `decline_reason` TEXT NULL DEFAULT NULL AFTER `proposed_datetime`;