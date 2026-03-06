package repo_mysql

import (
	"GoWebUser/models/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strconv"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrEmailNotExists = errors.New("email not exists")
	ErrIdNotExists    = errors.New("id not exists")
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &MySQLUserRepository{db: db}
}

// CreateUser 创建新用户
func (m *MySQLUserRepository) CreateUser(user *model.User) error {
	sqlStr := `insert into user(username,email,password) values(?,?,?)`
	_, err := m.db.Exec(sqlStr, user.Username, user.Email, user.Password)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// 转换为业务错误
			return ErrEmailExists
		}
		return fmt.Errorf("Exec failed,err:%w", err)
	}
	return nil
}

// GetByEmail 根据邮箱查找密码和id
func (m *MySQLUserRepository) GetByEmail(user *model.User) (*model.User, error) {
	row := m.db.QueryRow("select id,password from user where email=?", user.Email)
	u := &model.User{}
	err := row.Scan(&u.ID, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailNotExists
		}
		return nil, fmt.Errorf("Scan failed,err:%w", err)
	}
	return u, nil
}

// UpdateLastTime 更新用户最后一次登录的时间
func (m *MySQLUserRepository) UpdateLastTime(user *model.User) error {
	sqlStr := `update user set last_time = current_timestamp where id = ?`
	_, err := m.db.Exec(sqlStr, user.ID)
	if err != nil {
		return fmt.Errorf("Exec failed,err:%w", err)
	}
	return nil
}

// GetByUserID 根据用户id查找用户信息
func (m *MySQLUserRepository) GetByUserID(user *model.User) (*model.User, error) {
	row := m.db.QueryRow("select username,email,status,role,avatar,create_time,last_time from user where id=?", user.ID)
	u := &model.User{}
	err := row.Scan(&u.Username, &u.Email, &u.Status, &u.Role, &u.Avatar, &u.CreateTime, &u.LastTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrIdNotExists
		}
		return nil, fmt.Errorf("Scan failed,err:%w", err)
	}
	return u, nil
}

// GetLimitUsers 分页查询该页全部用户信息(可以带条件)
func (m *MySQLUserRepository) GetLimitUsers(page, size int, status, keyword string) ([]*model.User, int, error) {
	offest := (page - 1) * size
	baseSQL := "from user where 1=1 "
	args := []interface{}{}

	if status != "" {
		baseSQL += "and status = ? "
		stat, err := strconv.Atoi(status)
		if err != nil {
			return nil, 0, fmt.Errorf("Atoi failed,err:%w", err)
		}
		args = append(args, stat)
	}
	if keyword != "" {
		baseSQL += "and username like ? "
		args = append(args, "%"+keyword+"%")
	}

	//查询用户总数
	countSQL := "select count(*) " + baseSQL
	var num int
	row := m.db.QueryRow(countSQL, args...)
	err := row.Scan(&num)
	if err != nil {
		return nil, 0, fmt.Errorf("Scan failed,err:%w", err)
	}

	//查询用户信息
	dataSQL := "select id,username,email,status,role,avatar,create_time,last_time " + baseSQL + "limit ?,?"
	args = append(args, offest, size)
	rows, err := m.db.Query(dataSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Query failed,err:%w", err)
	}
	userList := []*model.User{}
	for rows.Next() {
		u := &model.User{}
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Status, &u.Role, &u.Avatar, &u.CreateTime, &u.LastTime)
		if err != nil {
			return nil, 0, fmt.Errorf("Scan failed,err:%w", err)
		}
		userList = append(userList, u)
	}
	return userList, num, nil
}

// DeleteUserByID 删除用户
func (m *MySQLUserRepository) DeleteUserByID(user *model.User) error {
	sqlStr := `delete from user where id = ?`
	_, err := m.db.Exec(sqlStr, user.ID)
	if err != nil {
		return fmt.Errorf("Exec failed,err:%w\n", err)
	}
	return nil
}
