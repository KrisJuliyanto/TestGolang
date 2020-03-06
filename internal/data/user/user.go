package user

import (
	"context"
	"log"

	"go-tutorial-2020/pkg/errors"
	"go-tutorial-2020/pkg/firebaseclient"

	userEntity "go-tutorial-2020/internal/entity/user"

	"cloud.google.com/go/firestore"
	"github.com/jmoiron/sqlx"
	"google.golang.org/api/iterator"
)

type (
	// Data ...
	Data struct {
		db   *sqlx.DB
		c    *firestore.Client
		stmt map[string]*sqlx.Stmt
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

const (
	getAllUsers  = "GetAllUsers"
	qGetAllUsers = "SELECT * FROM user_test"

	insertUser  = "InsertUser"
	qInsertUser = "INSERT INTO user_test VALUES (?,?,?,?,?,?)"

	getUserByName  = "GetUserByName"
	qGetUserByName = "SELECT * FROM user_test WHERE nama_lengkap LIKE ?"

	updateByNIP  = "UpdateByNIP"
	qUpdateByNIP = "UPDATE user_test SET nama_lengkap = ? , tanggal_lahir = ? , jabatan = ? , email = ? WHERE nip LIKE ?"

	getMaxNIP  = "GetMaxNIP"
	qGetMaxNIP = "SELECT MAX(CAST(RIGHT(nip,6)AS INT)) FROM user_test"

	deleteByNIP  = "DeleteByNIP"
	qDeleteByNIP = "DELETE FROM user_test WHERE nip = ?"
)

var (
	readStmt = []statement{
		{getAllUsers, qGetAllUsers},
		{insertUser, qInsertUser},
		{getUserByName, qGetUserByName},
		{updateByNIP, qUpdateByNIP},
		{getMaxNIP, qGetMaxNIP},
		{deleteByNIP, qDeleteByNIP},
	}
)

// New ...
func New(db *sqlx.DB, fc *firebaseclient.Client) Data {
	d := Data{
		db: db,
		c:  fc.Client,
	}
	d.initStmt()
	return d
}

func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

// GetAllUsers digunakan untuk mengambil semua data user
func (d Data) GetAllUsers(ctx context.Context) ([]userEntity.User, error) {
	var (
		user  userEntity.User
		users []userEntity.User
		err   error
	)

	// Query ke database
	rows, err := d.stmt[getAllUsers].QueryxContext(ctx)

	// Looping seluruh row data
	for rows.Next() {
		// Insert row data ke struct user
		if err := rows.StructScan(&user); err != nil {
			return users, errors.Wrap(err, "[DATA][GetAllUsers] ")
		}
		// Tambahkan struct user ke array user
		users = append(users, user)
	}
	// Return users array
	return users, err
}

// InsertUser untuk memasukkan data user
func (d Data) InsertUser(ctx context.Context, user userEntity.User) error {
	_, err := d.stmt[insertUser].ExecContext(ctx,
		user.ID,
		user.Nip,
		user.Nama,
		user.TanggalLahir,
		user.Jabatan,
		user.Email)
	return err
}

//UpdateUser ...
func (d Data) UpdateUser(ctx context.Context, user userEntity.User) error {
	_, err := d.stmt[updateByNIP].ExecContext(ctx,
		user.Nama,
		user.TanggalLahir,
		user.Jabatan,
		user.Email,
		user.Nip)
	return err
}

//GetUserByName ....
func (d Data) GetUserByName(ctx context.Context, userNama string) (userEntity.User, error) {
	var user userEntity.User
	err := d.stmt[getUserByName].QueryRowxContext(ctx, userNama).StructScan(&user)
	return user, err
}

//GetMaxNIP ...
func (d Data) GetMaxNIP(ctx context.Context) (int, error) {
	var nip int
	err := d.stmt[getMaxNIP].QueryRowxContext(ctx).Scan(&nip)
	log.Println(nip)
	return nip, err
}

//DeleteByNIP ...
func (d Data) DeleteByNIP(ctx context.Context, nip string) error {
	log.Println(nip)
	_, err := d.stmt[deleteByNIP].ExecContext(ctx, nip)
	return err
}

//ViewDataUserFirebase ...
func (d Data) ViewDataUserFirebase(ctx context.Context) ([]userEntity.User, error) {
	var (
		testFirebaseService = d.c.Collection("User")
		userList            []userEntity.User
		err                 error
	)
	iter := testFirebaseService.Documents(ctx)
	for {
		doc, err := iter.Next()
		var user userEntity.User
		if err == iterator.Done {
			break
		}
		err = doc.DataTo(&user)
		userList = append(userList, user)
	}
	return userList, err
}

//InsertUserFirebase ...
func (d Data) InsertUserFirebase(ctx context.Context, user userEntity.User) error {
	testFirebaseService := d.c.Collection("User")
	_, err := testFirebaseService.Doc(user.Nip).Set(ctx, user)
	log.Println("MASUK LOH")
	return err
}

//InsertManyFirebase ...
func (d Data) InsertManyFirebase(ctx context.Context, userList []userEntity.User) error {
	testFirebaseService := d.c.Collection("User")
	var err error
	for _, o := range userList {
		var user userEntity.User
		user.ID = o.ID
		user.Nip = o.Nip
		user.Nama = o.Nama
		user.TanggalLahir = o.TanggalLahir
		user.Jabatan = o.Jabatan
		user.Email = o.Email
		_, err = testFirebaseService.Doc(user.Nip).Set(ctx, user)
	}
	return err
}
