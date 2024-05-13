-- name: CreateUserAccountDetails :exec
INSERT INTO accounts (user_id, account_name, account_number, bank_name)
VALUES ($1, $2, $3, $4);

-- name: UpdateUserAccountDetails :exec
UPDATE accounts SET account_name = $1, account_number = $2, bank_name = $3 WHERE user_id = $4;

-- name: GetUserAccountDetails :one
SELECT * FROM accounts
WHERE user_id = $1;


-- name: SubtractTotalCoin :exec
DO $$
DECLARE
    user_id_val UUID := $1;
    amount DECIMAL := $2;
BEGIN
    UPDATE accounts
    SET total_coin = CASE 
                        WHEN total_coin >= amount THEN total_coin - amount
                        ELSE total_coin 
                    END
    WHERE user_id = user_id_val;

    IF NOT EXISTS (SELECT 1 FROM accounts WHERE user_id = user_id_val AND total_coin >= amount) THEN
        RAISE EXCEPTION 'Insufficient funds';
    END IF;
END $$;



-- CREATE OR REPLACE FUNCTION subtract_from_balance(user_id_param UUID, amount DECIMAL) RETURNS VOID AS $$
-- BEGIN
--     -- Check if the user has enough balance
--     IF (SELECT total_coin FROM accounts WHERE user_id = user_id_param) < amount THEN
--         RAISE EXCEPTION 'Insufficient funds';
--     ELSE
--         -- Perform the update if the balance is sufficient
--         UPDATE accounts
--         SET total_coin = total_coin - amount
--         WHERE user_id = user_id_param;
--     END IF;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- name: SubtractTotalCoin2 :exec
-- SELECT subtract_from_balance('99dc65d6-6bec-4f72-8df0-18949eea1ff8', 200);


-- name: IncreaseTotalCoin :one
UPDATE accounts
SET total_coin = total_coin + $1
WHERE user_id = $2
RETURNING total_coin;
