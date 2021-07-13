package crypto_app

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"my_projects/crypto/pkg/models"
	"my_projects/crypto/tools"
	"net/http"
	"strconv"
)

const (
	defaultCommission = float64(0.01)
	bcryptCost        = 11
)

type Crypto interface {
	Alive(ctx context.Context) (output models.AliveResponse, err error)
	Sign(ctx context.Context, input *models.RegisterRequest) (output models.RegisterResponse, err error)
	LogIn(ctx context.Context, input *models.LogInRequest) (output models.RegisterResponse, err error)
	GetWallets(ctx context.Context) (output []*models.WalletsResponse, err error)
	Transaction(ctx context.Context, input models.TransactionRequest) (success bool, err error)
	GetTransactions(ctx context.Context, perPage int, pageNum int) (response models.GetTransactionResponse, err error)
}

type crypto struct {
	db *pgx.Conn
}

func (r *crypto) Alive(ctx context.Context) (output models.AliveResponse, err error) {
	preID, err := strconv.Atoi(ctx.Value(models.CtxKey("id")).(string))
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении user_id из контекста",
			http.StatusInternalServerError)
	}
	output.UserID = int32(preID)
	output.Text = "service is okay"
	return
}

func (r *crypto) Sign(ctx context.Context, input *models.RegisterRequest) (output models.RegisterResponse, err error) {
	const (
		checkEmailExistence = `select exists(select * from user_data where email = $1);`
		dbRequestToAddData  = `insert  into user_data (name, last_name, email, pass_hash) values 
			($1, $2, $3, $4) returning id`
	)
	var (
		passHash    []byte
		emailExists bool
		userID      int32
	)

	if err = checkThePass(input.Pass); err != nil {
		return
	}
	if err = checkTheUserData(input.Name, input.LastName); err != nil {
		return
	}

	tx, err := r.db.Begin()
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании транзакции", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			er := tx.Rollback()
			if er != nil {
				log.Printf("error while rolling up the transaction: %v", er)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = tools.NewErrorMessage(err, "Ошибка при tx.Commit", http.StatusInternalServerError)
		}
	}()

	if input.Name == "" || input.LastName == "" || input.Pass == "" || input.Email == "" {
		err = tools.NewErrorMessage(errors.New("bad request"), "Какое то из полей пустое", http.StatusBadRequest)
		return
	}

	if err = tx.QueryRowEx(ctx, checkEmailExistence, nil, input.Email).Scan(&emailExists); err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при проверке на сущестование емейла", http.StatusInternalServerError)
		return
	}

	if emailExists {
		err = tools.NewErrorMessage(errors.New("this email is already registered"),
			"Данный емейл уже зарегестрирован", http.StatusBadRequest)
		return
	}

	if !isEmailValid(input.Email) {
		err = tools.NewErrorMessage(errors.New("bad email"),
			"Невалидный емейл", http.StatusBadRequest)
		return
	}

	if passHash, err = bcrypt.GenerateFromPassword([]byte(input.Pass), bcryptCost); err != nil {
		err = tools.NewErrorMessage(err, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	if err = tx.QueryRowEx(ctx, dbRequestToAddData, nil, input.Name, input.LastName, input.Email, string(passHash)).Scan(&userID); err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при сохранении данных", http.StatusInternalServerError)
		return
	}

	err = createDefaultWalletsWithDefaultBalance(ctx, tx, userID)
	if err != nil {
		return
	}

	output.AccessToken, err = generateToken(userID)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании токена", http.StatusInternalServerError)
	}
	return
}

func (r *crypto) LogIn(ctx context.Context, input *models.LogInRequest) (output models.RegisterResponse, err error) {
	const (
		getPassData = `select pass_hash,id from user_data where email = $1;`
	)

	var (
		passHash string
		userID   int32
	)

	if !isEmailValid(input.Email) {
		err = tools.NewErrorMessage(errors.New("bad email"),
			"Невалидный емейл", http.StatusBadRequest)
		return
	}

	if err = r.db.QueryRowEx(ctx, getPassData, nil, input.Email).Scan(&passHash, &userID); err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении данынх по емейлу",
			http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(input.Pass)); err != nil {
		err = tools.NewErrorMessage(err, "Некорректный пароль", http.StatusInternalServerError)
		return
	}

	output.AccessToken, err = generateToken(userID)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании токена", http.StatusInternalServerError)
	}
	return
}

func (r *crypto) GetWallets(ctx context.Context) (output []*models.WalletsResponse, err error) {
	const queryToGetWallets = `select address, name, balance from addresses as a 
    	left join salary s on a.salary_id = s.id
	where user_id = $1;`

	preID, err := strconv.Atoi(ctx.Value(models.CtxKey("id")).(string))
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении user_id из контекста",
			http.StatusInternalServerError)
	}
	userID := int32(preID)

	rows, err := r.db.QueryEx(ctx, queryToGetWallets, nil, userID)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении данных по кошелькам",
			http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		local := new(models.WalletsResponse)
		err = rows.Scan(
			&local.Address,
			&local.Salary,
			&local.Balance)
		if err != nil {
			err = tools.NewErrorMessage(err, "Ошибка при сканировании кошелька",
				http.StatusInternalServerError)
			break
		}
		output = append(output, local)
	}

	return
}

func (r *crypto) Transaction(ctx context.Context, input models.TransactionRequest) (success bool, err error) {
	const (
		queryToMakeTransaction = `select make_transaction($1,$2,$3,$4)`
	)

	if input.FromAddress == input.ToAddress {
		err = tools.NewErrorMessage(err, "Адрес не может быть одним и тем же", http.StatusBadRequest)
		return
	}

	tx, err := r.db.Begin()
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании транзакции", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			er := tx.Rollback()
			if er != nil {
				log.Printf("error while rolling up the transaction: %v", er)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = tools.NewErrorMessage(err, "Ошибка при tx.Commit", http.StatusInternalServerError)
		}
	}()

	preID, err := strconv.Atoi(ctx.Value(models.CtxKey("id")).(string))
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении user_id из контекста",
			http.StatusInternalServerError)
	}
	userID := int32(preID)

	err = checkTheAddressesBelongToPerson(ctx, tx, []int32{input.FromAddress, input.ToAddress}, userID)
	if err != nil {
		return
	}

	err = tx.QueryRowEx(ctx, queryToMakeTransaction, nil, input.FromAddress,
		input.ToAddress, input.Amount, defaultCommission).Scan(&success)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при переводе средств",
			http.StatusInternalServerError)
	}
	return
}

func (r *crypto) GetTransactions(ctx context.Context, perPage int, pageNum int) (response models.GetTransactionResponse, err error) {
	const query = `
		select a_from.address as from_address, a_to.address as to_address,amount_dollars as sum, commission, 
				cast(create_at as text) as date, successful as success from transactions as t
		    left join addresses a_from on a_from.id = t.from_address
		    left join addresses a_to on a_to.id = t.to_address
		where a_from.user_id = $1 or a_to.user_id = $1
		order by create_at desc`

	if pageNum < 0 || perPage <= 0 {
		err = tools.NewErrorMessage(errors.New("bad query params"), "Невалидные query параметры",
			http.StatusBadRequest)
		return
	}

	preID, err := strconv.Atoi(ctx.Value(models.CtxKey("id")).(string))
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении user_id из контекста",
			http.StatusInternalServerError)
	}
	userID := int32(preID)

	rows, err := r.db.QueryEx(ctx, query, nil, userID)
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при получении данных по транзакциям",
			http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		local := new(models.SingleTransaction)
		err = rows.Scan(
			&local.FromAddress,
			&local.ToAddress,
			&local.Sum,
			&local.Commission,
			&local.Date,
			&local.Success)
		if err != nil {
			err = tools.NewErrorMessage(err, "Ошибка при сканировании списка транзакций",
				http.StatusInternalServerError)
			break
		}
		response.Items = append(response.Items, local)
	}

	response.Meta.Total = int32(len(response.Items))
	response.Meta.PagesCount = response.Meta.Total / int32(perPage)
	if response.Meta.Total%int32(perPage) != 0 {
		response.Meta.PagesCount++
	}
	response.Meta.PageNum = int32(pageNum)

	if pageNum >= int(response.Meta.PagesCount) {
		err = tools.NewErrorMessageEncodeIntoWriter(errors.New("bad query params"), "Слишком большое значение page_num",
			http.StatusBadRequest)
		return models.GetTransactionResponse{}, err
	}

	var (
		firstVal int
		lastVal  int
	)

	firstVal = perPage * pageNum
	lastVal = firstVal + perPage
	if lastVal > int(response.Meta.Total) {
		lastVal = int(response.Meta.Total)
	}

	response.Items = response.Items[firstVal:lastVal]
	return
}

func NewCrypto(db *pgx.Conn) Crypto {
	return &crypto{
		db: db,
	}
}
