-- name: CreateBooking :exec
INSERT INTO bookings (id, first_name, last_name, gender, birthday, launch_pad_id, destination_id, launch_date,
                      created_at, updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10);

-- name: DeleteBooking :one
DELETE
FROM bookings
WHERE id = $1
RETURNING id;

-- name: GetBookingByID :one
SELECT id,
       first_name,
       last_name,
       gender,
       birthday,
       launch_pad_id,
       destination_id,
       launch_date,
       created_at,
       updated_at
FROM bookings
WHERE id = $1;

-- name: ListBookings :many
SELECT id,
       first_name,
       last_name,
       gender,
       birthday,
       launch_pad_id,
       destination_id,
       launch_date,
       created_at,
       updated_at
FROM bookings
WHERE launch_date = coalesce(sqlc.narg('launch_date'), launch_date)
  AND launch_pad_id = coalesce(sqlc.narg('launch_pad_id'), launch_pad_id)
  AND destination_id = coalesce(sqlc.narg('destination_id'), destination_id)
ORDER BY created_at DESC LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');
