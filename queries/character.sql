-- name: GetCharacter :one
SELECT c."ID", c.name, c."maxHP", c."HP", c.str, c.dex, c.con, c.int, c.wis, c.cha, c.gp, c.sp, c.cp, c.ep, c.pp, c."raceID", c."subraceID"
FROM "Character" c
WHERE "ID" = $1

-- name: GetCharacters :many
Select c."ID", c.name, c."maxHP", c."HP", c.str, c.dex, c.con, c.int, c.wis, c.cha, c.gp, c.sp, c.cp, c.ep, c.pp
FROM "Character"