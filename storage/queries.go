package storage

const createUser = `
INSERT INTO users (id, created_at, updated_at, name, type)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name) DO NOTHING
RETURNING id, created_at, updated_at, name, type
`

const getUserByID = `SELECT * FROM users WHERE id=$1`

const createTelegramUser = `
INSERT INTO telegramUsers(id, created_at, updated_at, name, chat_id, user_id, active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (chat_id) DO NOTHING
RETURNING id, created_at, updated_at, name, chat_id, user_id, active
`

const getTelegramUserByChatId = `SELECT * FROM telegramUsers WHERE chat_id=$1`

const updateTelegramUserActive = `
UPDATE telegramUsers
SET active=$1
WHERE chat_id=$2
`

const createSiteParse = `
INSERT INTO siteParse(id, created_at, updated_at, name, url_site, type, last_fetched_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url_site) DO NOTHING
RETURNING id, created_at, updated_at, name, url_site, type, last_fetched_at
`

const getAllSiteParse = `SELECT * FROM siteParse`

const getSiteParseById = `SELECT * FROM siteParse WHERE id=$1`

const getSiteParseByName = `SELECT * FROM siteParse WHERE name=$1`

const getSiteParseByType = `SELECT * FROM siteParse WHERE type=$1`

const createSiteParseFollow = `
INSERT INTO siteParseFollows(id, created_at, updated_at, user_id, site_parse_id, active)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id, site_parse_id) DO NOTHING
RETURNING id, created_at, updated_at, user_id, site_parse_id, active;
`

const getSiteParseFollowsByUserID = `SELECT * FROM siteParseFollows WHERE user_id=$1`

const updateSiteParseFollowsActive = `
UPDATE siteParseFollows
SET active = $1, updated_at = $2
WHERE user_id = $3 AND site_parse_id = $4;
`

const createNewsElement = `
INSERT INTO newsElements(id, created_at, updated_at, site_parse_id, title, news_date, url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url) DO NOTHING;
`
