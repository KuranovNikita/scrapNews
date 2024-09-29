package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrEmptyUser = errors.New("empty user")

var ErrEmptySiteParse = errors.New("empty siteParse")

var ErrEmptySiteParseFollow = errors.New("empty siteParseFollow")

type Storage struct {
	db *sql.DB
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Type      string
}

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Type      string
}

type TelegramUser struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ChatID    string
	UserID    uuid.UUID
	Active    bool
}

type CreateTelegramUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ChatID    string
	UserID    uuid.UUID
	Active    bool
}

type UpdateTelegramUserActiveParams struct {
	ChatID string
	Active bool
}

type SiteParse struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	UrlSite       string
	Type          string
	LastFetchedAt time.Time
}

type CreateSiteParseParams struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	UrlSite       string
	Type          string
	LastFetchedAt time.Time
}

type SiteParseFollows struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	SiteParseID uuid.UUID
	Active      bool
}

type CreateSiteParseParamsFollows struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	SiteParseID uuid.UUID
	Active      bool
}

type UpdateSiteParseActiveParamsFollows struct {
	UpdatedAt   time.Time
	UserID      uuid.UUID
	SiteParseID uuid.UUID
	Active      bool
}

type NewsElement struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SiteParseID uuid.UUID
	Title       string
	NewsDate    string
	Url         string
}

type CreateNewsElementParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SiteParseID uuid.UUID
	Title       string
	NewsDate    string
	Url         string
}

func New(driverName string, dataSourceName string) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	storage := &Storage{
		db: db,
	}

	return storage, nil
}

func (s *Storage) SaveUser(userParams CreateUserParams) (User, error) {
	var user User
	err := s.db.QueryRow(createUser,
		userParams.ID,
		userParams.CreatedAt,
		userParams.UpdatedAt,
		userParams.Name,
		userParams.Type,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Type,
	)
	if err == sql.ErrNoRows {
		return user, ErrEmptyUser
	}
	if err != nil {
		return user, fmt.Errorf("%w", err)
	}

	return user, nil
}

func (s *Storage) GetUserByID(ID string) (User, error) {
	var user User
	err := s.db.QueryRow(getUserByID, ID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Type)
	if err == sql.ErrNoRows {
		return user, ErrEmptyUser
	}
	if err != nil {
		return user, fmt.Errorf("%w", err)
	}

	return user, nil
}

func (s *Storage) SaveTelegramUser(telegramUserParams CreateTelegramUserParams) (TelegramUser, error) {
	var telegramUser TelegramUser
	err := s.db.QueryRow(createTelegramUser,
		telegramUserParams.ID,
		telegramUserParams.CreatedAt,
		telegramUserParams.UpdatedAt,
		telegramUserParams.Name,
		telegramUserParams.ChatID,
		telegramUserParams.UserID,
		telegramUserParams.Active,
	).Scan(
		&telegramUser.ID,
		&telegramUser.CreatedAt,
		&telegramUser.UpdatedAt,
		&telegramUser.Name,
		&telegramUser.ChatID,
		&telegramUser.UserID,
		&telegramUser.Active,
	)
	if err == sql.ErrNoRows {
		return telegramUser, ErrEmptyUser
	}
	if err != nil {
		return telegramUser, fmt.Errorf("%w", err)
	}

	return telegramUser, nil
}

func (s *Storage) GetTelegramUserByChatId(chatID string) (TelegramUser, error) {
	var telegramUser TelegramUser
	err := s.db.QueryRow(getTelegramUserByChatId, chatID).Scan(
		&telegramUser.ID,
		&telegramUser.CreatedAt,
		&telegramUser.UpdatedAt,
		&telegramUser.Name,
		&telegramUser.ChatID,
		&telegramUser.UserID,
		&telegramUser.Active,
	)
	if err == sql.ErrNoRows {
		return telegramUser, ErrEmptyUser
	}
	if err != nil {
		return telegramUser, fmt.Errorf("%w", err)
	}

	return telegramUser, nil
}

func (s *Storage) UpdateTelegramUserActive(updateTelegramUserActiveParams UpdateTelegramUserActiveParams) error {
	_, err := s.db.Exec(updateTelegramUserActive,
		updateTelegramUserActiveParams.Active,
		updateTelegramUserActiveParams.ChatID,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (s *Storage) CreateSiteParse(createSiteParseParams CreateSiteParseParams) (SiteParse, error) {
	var siteParse SiteParse
	err := s.db.QueryRow(createSiteParse,
		createSiteParseParams.ID,
		createSiteParseParams.CreatedAt,
		createSiteParseParams.UpdatedAt,
		createSiteParseParams.Name,
		createSiteParseParams.UrlSite,
		createSiteParseParams.Type,
		createSiteParseParams.LastFetchedAt,
	).Scan(
		&siteParse.ID,
		&siteParse.CreatedAt,
		&siteParse.UpdatedAt,
		&siteParse.Name,
		&siteParse.UrlSite,
		&siteParse.Type,
		&siteParse.LastFetchedAt,
	)
	if err != nil {
		return siteParse, fmt.Errorf("%w", err)
	}

	return siteParse, nil
}

func (s *Storage) GetAllSiteParses() ([]SiteParse, error) {
	var siteParses []SiteParse
	rows, err := s.db.Query(getAllSiteParse)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sp SiteParse
		err := rows.Scan(&sp.ID, &sp.CreatedAt, &sp.UpdatedAt, &sp.Name, &sp.UrlSite, &sp.Type, &sp.LastFetchedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		siteParses = append(siteParses, sp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllSiteParses rows iteration failed: %w", err)
	}

	return siteParses, nil
}

// func (s *Storage) GetSiteParseByType(Type string) (SiteParse, error) {
// 	var siteParse SiteParse
// 	err := s.db.QueryRow(getSiteParseById, Type).Scan(
// 		&siteParse.ID,
// 		&siteParse.CreatedAt,
// 		&siteParse.UpdatedAt,
// 		&siteParse.UrlSite,
// 		&siteParse.Name,
// 		&siteParse.Type,
// 		&siteParse.LastFetchedAt,
// 	)
// 	if err == sql.ErrNoRows {
// 		return siteParse, ErrEmptyUser
// 	}
// 	if err != nil {
// 		return siteParse, fmt.Errorf("%w", err)
// 	}

// 	return siteParse, nil
// }

func (s *Storage) GetSiteParseById(ID uuid.UUID) (SiteParse, error) {
	var siteParse SiteParse
	err := s.db.QueryRow(getSiteParseById, ID).Scan(
		&siteParse.ID,
		&siteParse.CreatedAt,
		&siteParse.UpdatedAt,
		&siteParse.UrlSite,
		&siteParse.Name,
		&siteParse.Type,
		&siteParse.LastFetchedAt,
	)
	if err == sql.ErrNoRows {
		return siteParse, ErrEmptyUser
	}
	if err != nil {
		return siteParse, fmt.Errorf("%w", err)
	}

	return siteParse, nil
}

func (s *Storage) GetSiteParseByName(name string) (SiteParse, error) {
	var siteParse SiteParse
	err := s.db.QueryRow(getSiteParseByName, name).Scan(
		&siteParse.ID,
		&siteParse.CreatedAt,
		&siteParse.UpdatedAt,
		&siteParse.UrlSite,
		&siteParse.Name,
		&siteParse.Type,
		&siteParse.LastFetchedAt,
	)
	if err == sql.ErrNoRows {
		return siteParse, ErrEmptyUser
	}
	if err != nil {
		return siteParse, fmt.Errorf("%w", err)
	}

	return siteParse, nil
}

func (s *Storage) CreateSiteParseFollows(createSiteParseParamsFollows CreateSiteParseParamsFollows) (SiteParseFollows, error) {
	var siteParseFollows SiteParseFollows
	err := s.db.QueryRow(createSiteParseFollow,
		createSiteParseParamsFollows.ID,
		createSiteParseParamsFollows.CreatedAt,
		createSiteParseParamsFollows.UpdatedAt,
		createSiteParseParamsFollows.UserID,
		createSiteParseParamsFollows.SiteParseID,
		createSiteParseParamsFollows.Active,
	).Scan(
		&siteParseFollows.ID,
		&siteParseFollows.CreatedAt,
		&siteParseFollows.UpdatedAt,
		&siteParseFollows.UserID,
		&siteParseFollows.SiteParseID,
		&siteParseFollows.Active,
	)
	if err != nil {
		return siteParseFollows, fmt.Errorf("%w", err)
	}

	return siteParseFollows, nil
}

func (s *Storage) UpdateSiteParseFollowsActive(updateSiteParseActiveParamsFollows UpdateSiteParseActiveParamsFollows) error {
	_, err := s.db.Exec(updateSiteParseFollowsActive,
		updateSiteParseActiveParamsFollows.Active,
		updateSiteParseActiveParamsFollows.UpdatedAt,
		updateSiteParseActiveParamsFollows.UserID,
		updateSiteParseActiveParamsFollows.SiteParseID,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (s *Storage) GetSiteParseFollowsByUserID(userID uuid.UUID) ([]SiteParseFollows, error) {
	var siteParseFollows []SiteParseFollows
	rows, err := s.db.Query(getSiteParseFollowsByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var spf SiteParseFollows
		err := rows.Scan(&spf.ID, &spf.CreatedAt, &spf.UpdatedAt, &spf.UserID, &spf.SiteParseID, &spf.Active)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		siteParseFollows = append(siteParseFollows, spf)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return siteParseFollows, nil
}

func (s *Storage) SaveNewsElement(createNewsElementsParams CreateNewsElementParams) (NewsElement, error) {
	var newsElement NewsElement
	err := s.db.QueryRow(createNewsElement,
		createNewsElementsParams.ID,
		createNewsElementsParams.CreatedAt,
		createNewsElementsParams.UpdatedAt,
		createNewsElementsParams.SiteParseID,
		createNewsElementsParams.Title,
		createNewsElementsParams.NewsDate,
		createNewsElementsParams.Url,
	).Scan(
		&newsElement.ID,
		&newsElement.CreatedAt,
		&newsElement.UpdatedAt,
		&newsElement.SiteParseID,
		&newsElement.Title,
		&newsElement.NewsDate,
		&newsElement.Url,
	)
	if err != nil {
		return newsElement, fmt.Errorf("%w", err)
	}

	return newsElement, nil
}
