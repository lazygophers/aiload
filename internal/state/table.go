package state

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/lrpc/middleware/xerror"
)

var (
	_db *db.Client

	Channel       *db.Model[aiload.ModelChannel]
	User          *db.Model[aiload.ModelUser]
	ChannelAccess *db.Model[aiload.ModelChannelAccess]
	ModelAlias    *db.Model[aiload.ModelModelAlias]
	UserToken     *db.Model[aiload.ModelUserToken]
	UserAccess    *db.Model[aiload.ModelUserAccess]
)

func ConnectDatabase() (err error) {
	log.Info("try init database")
	_db, err = db.New(State.Config.Db,
		&aiload.ModelChannel{},
		&aiload.ModelUser{},
		&aiload.ModelChannelAccess{},
		&aiload.ModelModelAlias{},
		&aiload.ModelUserToken{},
		&aiload.ModelUserAccess{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	Channel = db.NewModel[aiload.ModelChannel](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ChannelNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ChannelDuplicateKey)))
	User = db.NewModel[aiload.ModelUser](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserDuplicateKey)))
	ChannelAccess = db.NewModel[aiload.ModelChannelAccess](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ChannelAccessNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ChannelAccessDuplicateKey)))
	ModelAlias = db.NewModel[aiload.ModelModelAlias](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ModelAliasNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ModelAliasDuplicateKey)))
	UserToken = db.NewModel[aiload.ModelUserToken](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserTokenNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserTokenDuplicateKey)))
	UserAccess = db.NewModel[aiload.ModelUserAccess](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserAccessNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserAccessDuplicateKey)))

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
