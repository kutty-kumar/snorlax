package entity

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kutty-kumar/charminder/pkg"
	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/ho_oh/core_v1"
	"github.com/kutty-kumar/ho_oh/user_service_v1"
)

type User struct {
	charminder.BaseDomain
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Gender      core_v1.Gender
	Email       string
	Password    string
}

func (r *User) ToJson() (string, error) {
	rBytes, err := json.Marshal(*r)
	if err != nil {
		return "{}", err
	}
	return string(rBytes), nil
}

func (r *User) MarshalBinary() ([]byte, error) {
	var rBytes bytes.Buffer
	enc := gob.NewEncoder(&rBytes)
	err := enc.Encode(*r)
	return rBytes.Bytes(), err
}

func (r *User) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(r); err != nil {
		return err
	}
	return nil
}

func (r *User) String() string {
	return fmt.Sprintf("{}")
}

func (r *User) GetName() pkg.DomainName {
	return "users"
}

func (r *User) ToDto() interface{} {
	return user_service_v1.UserDto{
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		DateOfBirth: r.DateOfBirth.Format("2006-01-02"),
		Gender:      r.Gender,
		UserId:      r.ExternalId,
		Status:      core_v1.Status(r.Status),
		Email:       r.Email,
	}
}

func (r *User) FillProperties(dto interface{}) charminder.Base {
	userDto := dto.(user_service_v1.UserDto)
	r.FirstName = userDto.FirstName
	r.LastName = userDto.LastName
	// err here is ignored as correctness is guaranteed
	parsedTime, _ := time.Parse("2006-01-02", userDto.DateOfBirth)
	r.DateOfBirth = parsedTime
	r.Gender = userDto.Gender
	r.Status = int(userDto.Status)
	r.Email = userDto.Email
	r.Password = userDto.Password
	return r
}

func (r *User) Merge(other interface{}) {
	updatableUser := other.(*User)
	if !updatableUser.DateOfBirth.IsZero() {
		r.DateOfBirth = updatableUser.DateOfBirth
	}
	if updatableUser.FirstName != "" {
		r.FirstName = updatableUser.FirstName
	}
	if updatableUser.LastName != "" {
		r.LastName = updatableUser.LastName
	}
	if updatableUser.Password != "" {
		r.Password = updatableUser.Password
	}
}

func (r *User) FromSqlRow(rows *sql.Rows) (charminder.Base, error) {
	err := rows.Scan(&r.Id, &r.ExternalId, &r.FirstName, &r.LastName, &r.DateOfBirth, &r.LastName, &r.Gender, &r.Status)
	return r, err
}

func (r *User) SetExternalId(externalId string) {
	r.ExternalId = externalId
}
