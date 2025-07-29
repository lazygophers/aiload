package state

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/db"

)

var (
	_db *db.Client

	UserToken     *db.Model[aiload.ModelUserToken]
	Channel       *db.Model[aiload.ModelChannel]
	ChannelAccess *db.Model[aiload.ModelChannelAccess]
	UserAccess    *db.Model[aiload.ModelUserAccess]
	User          *db.Model[aiload.ModelUser]
	ModelAlias    *db.Model[aiload.ModelModelAlias]
)

func ConnectDatabase() (err error) {
	log.Info("try init database")
	_db, err = db.New(State.Config.Db,
		&aiload.ModelUserToken{},
		&aiload.ModelChannel{},
		&aiload.ModelChannelAccess{},
		&aiload.ModelUserAccess{},
		&aiload.ModelUser{},
		&aiload.ModelModelAlias{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	UserToken = db.NewModel[aiload.ModelUserToken](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserTokenNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserTokenDuplicateKey)))
	Channel = db.NewModel[aiload.ModelChannel](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ChannelNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ChannelDuplicateKey)))
	ChannelAccess = db.NewModel[aiload.ModelChannelAccess](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ChannelAccessNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ChannelAccessDuplicateKey)))
	UserAccess = db.NewModel[aiload.ModelUserAccess](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserAccessNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserAccessDuplicateKey)))
	User = db.NewModel[aiload.ModelUser](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserDuplicateKey)))
	ModelAlias = db.NewModel[aiload.ModelModelAlias](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ModelAliasNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ModelAliasDuplicateKey)))

	log.Info("connect database successfully")

	return nil
}

func Db() *db.Client {
	return _db
}

func NewScoop() *db.Scoop {
	return _db.NewScoop()
}

func Begin() *db.Scoop {
	return NewScoop().Begin()
}

func CommitOrRollback(logic func(tx *db.Scoop) error) error {
	return NewScoop().CommitOrRollback(NewScoop().Begin(), logic)
}
