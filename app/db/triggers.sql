CREATE OR REPLACE FUNCTION reset_customer_seq() RETURNS trigger AS $$
BEGIN
    -- Check if the customer table is empty.
    IF (SELECT COUNT(*) FROM customer) = 0 THEN
        EXECUTE 'ALTER SEQUENCE customer_customer_id_seq RESTART WITH 1';
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reset_customer_seq_trigger
AFTER DELETE ON customer
FOR EACH STATEMENT
EXECUTE FUNCTION reset_customer_seq();




CREATE OR REPLACE FUNCTION reset_vendor_seq() RETURNS trigger AS $$
BEGIN
    -- Check if the vendor table is empty.
    IF (SELECT COUNT(*) FROM vendor) = 0 THEN
        EXECUTE 'ALTER SEQUENCE vendor_vendor_id_seq RESTART WITH 1';
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reset_vendor_seq_trigger
AFTER DELETE ON vendor
FOR EACH STATEMENT
EXECUTE FUNCTION reset_vendor_seq();
