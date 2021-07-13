package crypto_app

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx"
	"github.com/lib/pq"
	"math/rand"
	"my_projects/crypto/pkg/models"
	"my_projects/crypto/tools"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	defaultBalance = float64(100)
	btc            = 1
	eth            = 2
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	preValid := emailRegex.MatchString(email)
	if !preValid {
		return preValid
	}

	some := strings.Split(email, "@")
	if len(some) != 2 {
		return false
	}

	_, err := net.LookupIP(some[1])
	return err == nil
}

func checkThePass(pass string) (err error) {
	if len(pass) != len([]rune(pass)) {
		err = tools.NewErrorMessage(errors.New("bad pass"),
			"Некорректный пароль",
			http.StatusBadRequest)
		return
	}

	if len(pass) < 8 || len(pass) > 50 {
		err = tools.NewErrorMessage(errors.New("bad pass len"),
			"Длина пароля должна быть от 8 до 50 символов", http.StatusBadRequest)
	}
	return
}

func checkTheUserData(firstName, lastName string) (err error) {
	fieldsArr := []string{firstName, lastName}

	for i := range fieldsArr {
		fieldRune := []rune(fieldsArr[i])
		for i := range fieldRune {
			if !unicode.IsLetter(fieldRune[i]) {
				err = tools.NewErrorMessage(errors.New("bad field"),
					"Следует использовать только символы", http.StatusBadRequest)
				return
			}
		}
	}
	return
}

func generateToken(userID int32) (response string, err error) {
	claims := models.ClaimWithID{
		ID: strconv.Itoa(int(userID)),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	response, err = token.SignedString(models.JwtSigningKey)
	return
}

func createDefaultWalletsWithDefaultBalance(ctx context.Context, tx *pgx.Tx, userID int32) (err error) {
	const (
		queryToAddNewWallet = `insert into addresses (address, user_id, salary_id, balance) values ($1,$2,$3,$4);`
	)

	address := randNumberString(20)
	rows, err := tx.QueryEx(ctx, queryToAddNewWallet, nil, address, userID, btc, defaultBalance)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании кошелька",
			http.StatusInternalServerError)
		return
	}
	rows.Close()

	address = randNumberString(20)
	rows, err = tx.QueryEx(ctx, queryToAddNewWallet, nil, address, userID, eth, defaultBalance)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании кошелька",
			http.StatusInternalServerError)
	}
	rows.Close()
	return
}

func randNumberString(n int) string {
	var addressSymbols = "1234567890"
	b := make([]uint8, n)
	for i := range b {
		b[i] = addressSymbols[rand.Intn(len(addressSymbols))]
	}
	return string(b)
}

func checkTheAddressesBelongToPerson(ctx context.Context, tx *pgx.Tx, addresses []int32, userID int32) (err error) {
	var users []int32
	const query = `select distinct user_id from addresses where id = any($1);`

	rows, err := tx.QueryEx(ctx, query, nil, pq.Array(addresses))
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении данных о адресах",
			http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var localUser int32
		err = rows.Scan(&localUser)
		if err != nil {
			err = tools.NewErrorMessage(err, "Ошибка при сканировании юзера при проверке",
				http.StatusInternalServerError)
			break
		}
		users = append(users, localUser)
	}

	if len(users) == 1 && users[0] == userID {
		return nil
	}
	err = tools.NewErrorMessage(errors.New("bad addresses"), "В данном наборе адресов есть адреса принадлежащие нескольким пользователям",
		http.StatusBadRequest)
	return
}
