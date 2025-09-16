CREATE TABLE `payment_methods` (
   `id` INT NOT NULL AUTO_INCREMENT,
   `user_id` INT NOT NULL,
   `midtrans_token` VARCHAR(255) NOT NULL,
    `card_type` VARCHAR(50) NULL,
    `masked_card` VARCHAR(50) NOT NULL,
    `is_default` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `midtrans_token_UNIQUE` (`midtrans_token` ASC),
    INDEX `fk_payment_methods_user_id_idx` (`user_id` ASC),
    CONSTRAINT `fk_payment_methods_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
                                                              ON DELETE CASCADE
    ) ENGINE=InnoDB;

-- Vehicle Images table
CREATE TABLE `vehicle_images` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `vehicle_id` INT NOT NULL,
  `image_url` VARCHAR(255) NOT NULL,
    `is_primary` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_vehicle_images_vehicle_id_idx` (`vehicle_id` ASC),
    CONSTRAINT `fk_vehicle_images_vehicle_id`
    FOREIGN KEY (`vehicle_id`)
    REFERENCES `vehicles` (`id`)
                                                              ON DELETE CASCADE
    ) ENGINE=InnoDB;

-- Bookings table to track customer interest
CREATE TABLE `bookings` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `vehicle_id` INT NOT NULL,
    `booking_date` TIMESTAMP NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_bookings_user_id_idx` (`user_id` ASC),
    INDEX `fk_bookings_vehicle_id_idx` (`vehicle_id` ASC),
    CONSTRAINT `fk_bookings_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
                                                              ON DELETE RESTRICT,
    CONSTRAINT `fk_bookings_vehicle_id`
    FOREIGN KEY (`vehicle_id`)
    REFERENCES `vehicles` (`id`)
                                                              ON DELETE RESTRICT
    ) ENGINE=InnoDB;

-- Schedules table for appointments
CREATE TABLE `schedules` (
     `id` INT NOT NULL AUTO_INCREMENT,
     `booking_id` INT NOT NULL,
     `user_id` INT NOT NULL, -- The staff member (user) assigned
     `appointment_datetime` TIMESTAMP NOT NULL,
     `notes` TEXT NULL,
     `status` VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `booking_id_UNIQUE` (`booking_id` ASC),
    INDEX `fk_schedules_user_id_idx` (`user_id` ASC),
    CONSTRAINT `fk_schedules_booking_id`
    FOREIGN KEY (`booking_id`)
    REFERENCES `bookings` (`id`)
                                                              ON DELETE CASCADE,
    CONSTRAINT `fk_schedules_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
                                                              ON DELETE RESTRICT
    ) ENGINE=InnoDB;

-- Agreements table to formalize a deal
CREATE TABLE `agreements` (
      `id` INT NOT NULL AUTO_INCREMENT,
      `booking_id` INT NOT NULL,
      `agreement_date` TIMESTAMP NOT NULL,
      `final_price` DECIMAL(15, 2) NOT NULL,
    `terms` TEXT NULL,
    `signed_by_user` BOOLEAN NOT NULL DEFAULT FALSE,
    `signed_by_staff` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `booking_id_UNIQUE` (`booking_id` ASC),
    CONSTRAINT `fk_agreements_booking_id`
    FOREIGN KEY (`booking_id`)
    REFERENCES `bookings` (`id`)
      ON DELETE RESTRICT
    ) ENGINE=InnoDB;

-- Payments table, designed for Midtrans
CREATE TABLE `payments` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `agreement_id` INT NOT NULL,
    `amount` DECIMAL(15, 2) NOT NULL,
    `payment_method` VARCHAR(50) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `midtrans_transaction_id` VARCHAR(255) NULL,
    `payment_url` VARCHAR(255) NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `midtrans_transaction_id_UNIQUE` (`midtrans_transaction_id` ASC),
    INDEX `fk_payments_agreement_id_idx` (`agreement_id` ASC),
    CONSTRAINT `fk_payments_agreement_id`
    FOREIGN KEY (`agreement_id`)
    REFERENCES `agreements` (`id`)
        ON DELETE RESTRICT
    ) ENGINE=InnoDB;

-- Installments table with penalty tracking
CREATE TABLE `installments` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `payment_id` INT NOT NULL,
    `due_date` DATE NOT NULL,
    `amount_due` DECIMAL(15, 2) NOT NULL,
    `penalty_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `total_due` DECIMAL(15, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `paid_date` DATE NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_installments_payment_id_idx` (`payment_id` ASC),
    CONSTRAINT `fk_installments_payment_id`
    FOREIGN KEY (`payment_id`)
    REFERENCES `payments` (`id`)
    ON DELETE CASCADE
    ) ENGINE=InnoDB;