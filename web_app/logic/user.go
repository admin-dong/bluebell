package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	snowflake "web_app/pkg"
	"web_app/pkg/jwt"
)

func SignUp(p *models.ParamSignUP) (err error) {
	//1.判断与用户存在不存在
	if err = mysql.CheckUserExit(p.Username); err != nil {
		//数据库查询出错
		return err
	}
	//2生成uid
	userID := snowflake.GenID()
	//构造一个user示例
	user := &models.User{
		userID,
		p.Username,
		p.Password,
		"_",
	}
	//3保存进数据库

	return mysql.InsertUser(user)

}

//登录。
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针  就能拿到user。userid
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成jwt
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
