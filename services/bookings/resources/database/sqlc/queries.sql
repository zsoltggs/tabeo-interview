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

-- name: DeleteBooking :exec
DELETE FROM bookings WHERE id = $1;

-- name: GetBookingByID :one
SELECT
    id, first_name, last_name, gender, birthday, launch_pad_id, destination_id, launch_date, created_at, updated_at
FROM
    bookings
WHERE
    id = $1;

-- name: ListBookings :many
SELECT
    id, first_name, last_name, gender, birthday, launch_pad_id, destination_id, launch_date, created_at, updated_at
FROM
    bookings
WHERE
    ($1::timestamptz IS NULL OR launch_date = $1)
  AND ($2::varchar IS NULL OR launch_pad_id = $2)
  AND ($3::varchar IS NULL OR destination_id = $3)
ORDER BY
    created_at DESC
    LIMIT $4 OFFSET $5;
