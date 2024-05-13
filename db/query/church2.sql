-- -- Change Church Membership
-- -- :user_id is the user's ID, :church_id is the church's ID
-- -- :join_date is the timestamp of the join date
-- -- This query assumes the existence of user_church_membership table
-- INSERT INTO user_church_membership (user_id, church_id, join_date)
-- VALUES (:user_id, :church_id, :join_date)
-- ON CONFLICT (user_id) DO UPDATE SET church_id = EXCLUDED.church_id, join_date = EXCLUDED.join_date;

-- -- Change Denomination Membership
-- -- :user_id is the user's ID, :denomination_id is the denomination's ID
-- -- :join_date is the timestamp of the join date
-- -- This query assumes the existence of user_denomination_membership table
-- INSERT INTO user_denomination_membership (user_id, denomination_id, join_date)
-- VALUES (:user_id, :denomination_id, :join_date)
-- ON CONFLICT (user_id) DO UPDATE SET denomination_id = EXCLUDED.denomination_id, join_date = EXCLUDED.join_date;

-- -- Search Churches
-- -- :user_id is the user's ID, :query is the search query
-- -- :latitude and :longitude are optional for nearby search
-- SELECT c.id, c.denomination_id, c.name, c.vicar, c.location
-- FROM churches c
-- LEFT JOIN user_denomination_membership udm ON c.denomination_id = udm.denomination_id
-- WHERE c.name ILIKE :query AND (udm.user_id = :user_id OR udm.user_id IS NULL)
-- ORDER BY udm.user_id IS NOT NULL DESC;

-- -- Select Denomination
-- -- :user_id is the user's ID, :denomination_id is the denomination's ID
-- -- This query assumes the existence of user_denomination_membership table
-- INSERT INTO user_denomination_membership (user_id, denomination_id, join_date)
-- VALUES (:user_id, :denomination_id, NOW())
-- ON CONFLICT (user_id) DO UPDATE SET denomination_id = EXCLUDED.denomination_id, join_date = EXCLUDED.join_date;

-- -- Create Church
-- -- :name is the church's name, :vicar is the vicar's name
-- -- :latitude and :longitude are the church's location
-- -- :denomination_id is the denomination's ID
-- INSERT INTO churches (name, vicar, location, denomination_id)
-- VALUES (:name, :vicar, ST_MakePoint(:longitude, :latitude), :denomination_id)
-- RETURNING id;

-- -- Search Nearby Churches
-- -- :user_id is the user's ID, :query is the search query
-- -- :latitude and :longitude are the user's location
-- SELECT c.id, c.denomination_id, c.name, c.vicar, c.location
-- FROM churches c
-- LEFT JOIN user_denomination_membership udm ON c.denomination_id = udm.denomination_id
-- WHERE ST_Distance(c.location, ST_MakePoint(:longitude, :latitude)) < 5000  -- Assuming 5000 meters radius
-- ORDER BY ST_Distance(c.location, ST_MakePoint(:longitude, :latitude));
