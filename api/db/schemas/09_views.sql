CREATE
OR REPLACE VIEW user_view AS
SELECT
  *
FROM
  "user"
WHERE
  is_deleted = FALSE;