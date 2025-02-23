-- Removed admin Table
-- Removed chat Table
-- Removed support Table

-- Updated customer Table (No changes needed here)
CREATE TABLE IF NOT EXISTS customer (
    customer_id    SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    phone         VARCHAR(20),
    email         VARCHAR(100) NOT NULL UNIQUE,
    address       VARCHAR(255)
);

-- Updated vendor Table (No changes needed here)
CREATE TABLE IF NOT EXISTS vendor (
    vendor_id       SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    email          VARCHAR(100) NOT NULL UNIQUE,
    phone          VARCHAR(20),
    hotel_name     VARCHAR(255),
    address        VARCHAR(255)
);

-- Updated room Table (No changes needed here)
CREATE TABLE IF NOT EXISTS room (
    room_id       SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    location      VARCHAR(255),
    availability  BOOLEAN DEFAULT TRUE,
    price         NUMERIC(10,2) NOT NULL,
    room_type     VARCHAR(50),
    average_rating NUMERIC(3,2) DEFAULT 0.0,
    amenities     TEXT DEFAULT '',  -- Changed from TEXT[] to TEXT
    vendor_id     INT NOT NULL,
    CONSTRAINT fk_room_vendor FOREIGN KEY (vendor_id) REFERENCES vendor(vendor_id) ON DELETE CASCADE
);

-- Updated booking Table (No changes needed here)
CREATE TABLE IF NOT EXISTS booking (
    booking_id     SERIAL PRIMARY KEY,
    booking_date   DATE NOT NULL DEFAULT CURRENT_DATE,
    checkin_date   DATE NOT NULL,
    checkout_date  DATE NOT NULL,
    payment_status VARCHAR(50) CHECK (payment_status IN ('Pending', 'Paid', 'Failed')),
    room_id        INT NOT NULL,
    customer_id    INT NOT NULL,
    CONSTRAINT fk_booking_room FOREIGN KEY (room_id) REFERENCES room(room_id) ON DELETE CASCADE,
    CONSTRAINT fk_booking_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id) ON DELETE CASCADE
);

-- Updated payment Table (No changes needed here)
CREATE TABLE IF NOT EXISTS payment (
    payment_id      SERIAL PRIMARY KEY,
    payment_method  VARCHAR(50),
    payment_status  VARCHAR(50) CHECK (payment_status IN ('Pending','Completed', 'Failed')),
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount         NUMERIC(10,2) NOT NULL,
    booking_id     INT NOT NULL,
    CONSTRAINT fk_payment_booking FOREIGN KEY (booking_id) REFERENCES booking(booking_id) ON DELETE CASCADE
);

-- Updated review Table (No changes needed here)
CREATE TABLE IF NOT EXISTS review (
    review_id    SERIAL PRIMARY KEY,
    comment      TEXT,
    rating       INT CHECK (rating >= 1 AND rating <= 5),
    review_date  DATE NOT NULL DEFAULT CURRENT_DATE,
    booking_id   INT NOT NULL,
    customer_id  INT NOT NULL,
    room_id      INT NOT NULL,
    CONSTRAINT fk_review_booking FOREIGN KEY (booking_id) REFERENCES booking(booking_id) ON DELETE CASCADE,
    CONSTRAINT fk_review_customer FOREIGN KEY (customer_id) REFERENCES customer(customer_id) ON DELETE CASCADE,
    CONSTRAINT fk_review_room FOREIGN KEY (room_id) REFERENCES room(room_id) ON DELETE CASCADE
);
