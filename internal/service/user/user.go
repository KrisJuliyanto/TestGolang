package skeleton

import (
	"context"
	"fmt"
	"log"
	"strings"

	userEntity "go-tutorial-2020/internal/entity/user"
	"go-tutorial-2020/pkg/errors"
	"go-tutorial-2020/pkg/kafka"
)

// UserData ...
type UserData interface {
	GetAllUsers(ctx context.Context) ([]userEntity.User, error)
	InsertUser(ctx context.Context, user userEntity.User) error
	GetUserByName(ctx context.Context, userNama string) (userEntity.User, error)
	UpdateUser(ctx context.Context, user userEntity.User) error
	GetMaxNIP(ctx context.Context) (int, error)
	DeleteByNIP(ctx context.Context, nip string) error
	ViewDataUserFirebase(ctx context.Context) ([]userEntity.User, error)
	InsertUserFirebase(ctx context.Context, user userEntity.User) error
	InsertManyFirebase(ctx context.Context, userList []userEntity.User) error
}

// Service ...
type Service struct {
	userData UserData
	kafkaSvc *kafka.Kafka
}

// New ...
func New(userData UserData, kafkaSvc *kafka.Kafka) Service {
	return Service{
		userData: userData,
		kafkaSvc: kafkaSvc,
	}
}

// GetAllUsers ...
func (s Service) GetAllUsers(ctx context.Context) ([]userEntity.User, error) {
	// Panggil method GetAllUsers di data layer user
	users, err := s.userData.GetAllUsers(ctx)
	// Error handling
	if err != nil {
		return users, errors.Wrap(err, "[SERVICE][GetAllUsers]")
	}
	// Return users array
	return users, err
}

// InsertUser ...
func (s Service) InsertUser(ctx context.Context, user userEntity.User) error {
	var (
		userValidasi userEntity.User
		err          error
		maxNip       int
	)
	userValidasi, err = s.userData.GetUserByName(ctx, user.Nama)
	if strings.Contains(userValidasi.Nama, user.Nama) {
		return errors.Wrap(errors.New("data already exists"), "[SERVICE][InsertUser]")
	}
	maxNip, err = s.userData.GetMaxNIP(ctx)
	user.Nip = "P" + fmt.Sprintf("%06d", (maxNip+1))
	log.Println(user.Nip)
	err = s.userData.InsertUser(ctx, user)
	return err
}

//UpdateUser ...
func (s Service) UpdateUser(ctx context.Context, user userEntity.User) error {
	var (
		//userValidasi userEntity.User
		err error
	)
	// userValidasi, err = s.userData.GetUserByName(ctx, user.Nama)
	// if userValidasi.Nama != user.Nama {
	// 	return errors.Wrap(errors.New("data already exists"), "[SERVICE][InsertUser]")
	// }
	err = s.userData.UpdateUser(ctx, user)
	return err
}

// GetUserByName ...
func (s Service) GetUserByName(ctx context.Context, userNama string) (userEntity.User, error) {
	result, err := s.userData.GetUserByName(ctx, userNama)
	return result, err
}

//DeleteByNIP ...
func (s Service) DeleteByNIP(ctx context.Context, nip string) error {
	err := s.userData.DeleteByNIP(ctx, nip)
	return err
}

//ViewDataUserFirebase ...
func (s Service) ViewDataUserFirebase(ctx context.Context) ([]userEntity.User, error) {
	userList, err := s.userData.ViewDataUserFirebase(ctx)
	return userList, err
}

//InsertUserFirebase ..
func (s Service) InsertUserFirebase(ctx context.Context, user userEntity.User) error {
	err := s.userData.InsertUserFirebase(ctx, user)
	return err
}

//InsertManyFirebase ...
func (s Service) InsertManyFirebase(ctx context.Context, userList []userEntity.User) error {
	err := s.userData.InsertManyFirebase(ctx, userList)
	if err != nil {
		log.Println(err)
	}
	return err
}

//PublishUser ...
func (s Service) PublishUser(user userEntity.User) error {
	err := s.kafkaSvc.SendMessageJSON("New_User", user)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][PublishUser]")
	}
	return err
}
