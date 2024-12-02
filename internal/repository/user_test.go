package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/pkg/geo"
	hash "github.com/a-berahman/dating-app/pkg/hash"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}
func TestCreateUser(t *testing.T) {
	db, mock, err := NewMock()
	assert.NoError(t, err)
	repo := repo{db}

	testCases := []struct {
		name        string
		email       string
		password    string
		personName  string
		gender      constant.UserGender
		dateOfBirth time.Time
		lat         float64
		lng         float64
		expectError bool
	}{
		{
			name:        "Successful Creation",
			email:       "test@example.com",
			password:    "123",
			personName:  "test name",
			gender:      constant.UserGenderMale,
			dateOfBirth: time.Now(),
			lat:         34.0522,
			lng:         -118.2437,
			expectError: false,
		},
		{
			name:        "Failed Creation",
			email:       "fail@example.com",
			password:    "123",
			personName:  "fail test name",
			gender:      constant.UserGenderFemale,
			dateOfBirth: time.Now(),
			lat:         0.0,
			lng:         0.0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			location, _ := geo.GeoEncode(tc.lat, tc.lng)

			mock.ExpectBegin()
			if !tc.expectError {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password","name","gender","date_of_birth","location") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tc.email, sqlmock.AnyArg(), tc.personName, tc.gender, sqlmock.AnyArg(), location).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			} else {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			}

			_, err := repo.Create(context.Background(), tc.email, tc.password, tc.personName, tc.gender, tc.dateOfBirth, tc.lat, tc.lng)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			}

		})
	}
}

func TestAuthenticate(t *testing.T) {
	db, mock, err := NewMock()
	assert.NoError(t, err)
	repo := repo{db}

	hashedPassword, _ := hash.Generate([]byte("passssssword"), bcrypt.DefaultCost)

	testCases := []struct {
		name          string
		email         string
		password      string
		setupMock     func()
		expectSuccess bool
		expectError   bool
	}{
		{
			name:     "Successful Authentication",
			email:    "testmail@example.com",
			password: "passssssword",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "testmail@example.com", string(hashedPassword))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("testmail@example.com", 1). // Adding correct placeholder for LIMIT
					WillReturnRows(rows)
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name:     "User Not Found",
			email:    "testmail@example.com",
			password: "passssssword",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("testmail@example.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectSuccess: false,
			expectError:   true,
		},
		{
			name:     "Invalid Password",
			email:    "testmail@example.com",
			password: "wrongPassssssword",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "testmail@example.com", string(hashedPassword))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("testmail@example.com", 1).
					WillReturnRows(rows)
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			user, err := repo.Authenticate(context.Background(), tc.email, tc.password)
			if tc.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFindByEmail(t *testing.T) {
	db, mock, err := NewMock()
	assert.NoError(t, err)
	repo := repo{db}

	testCases := []struct {
		name        string
		email       string
		setupMock   func()
		expectUser  *User
		expectError bool
	}{
		{
			name:  "Successful Find",
			email: "ahmad@test.com",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "name", "gender", "date_of_birth", "location"}).
					AddRow(1, "ahmad@test.com", "passsssssswwwooord", "ahmad", "MALE", time.Now(), "location")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs("ahmad@test.com", 1).
					WillReturnRows(rows)
			},
			expectUser: &User{
				Model:       gorm.Model{ID: 1},
				Email:       "ahmad@test.com",
				Password:    "passsssssswwwooord",
				Name:        "ahmad",
				Gender:      "MALE",
				DateOfBirth: time.Now(),
				Location:    "location",
			},
			expectError: false,
		},
		{
			name:  "User Not Found",
			email: "test",
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."email" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
					WithArgs(999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectUser:  nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			user, err := repo.FindByEmail(context.Background(), tc.email)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectUser.ID, user.ID)
				assert.Equal(t, tc.expectUser.Email, user.Email)
				assert.Equal(t, tc.expectUser.Name, user.Name)
			}

		})
	}
}
