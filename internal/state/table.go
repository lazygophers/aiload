package state

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/lrpc/middleware/xerror"
)

var (
	_db *db.Client

	Api           *db.Model[aiload.ModelApi]
	ApiAccess     *db.Model[aiload.ModelApiAccess]
	ModelAliasMap *db.Model[aiload.ModelModelAliasMap]
	User          *db.Model[aiload.ModelUser]
	UserToken     *db.Model[aiload.ModelUserToken]
	UserAccess    *db.Model[aiload.ModelUserAccess]
)

func ConnectDatabase() (err error) {
	log.Info("try init database")
	_db, err = db.New(State.Config.Db,
		&aiload.ModelApi{},
		&aiload.ModelApiAccess{},
		&aiload.ModelModelAliasMap{},
		&aiload.ModelUser{},
		&aiload.ModelUserToken{},
		&aiload.ModelUserAccess{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	Api = db.NewModel[aiload.ModelApi](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ApiNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ApiDuplicateKey)))
	ApiAccess = db.NewModel[aiload.ModelApiAccess](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ApiAccessNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ApiAccessDuplicateKey)))
	ModelAliasMap = db.NewModel[aiload.ModelModelAliasMap](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_ModelAliasMapNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_ModelAliasMapDuplicateKey)))
	User = db.NewModel[aiload.ModelUser](Db()).
		SetNotFound(xerror.NewError(int32(aiload.ErrCode_UserNotFound))).
		SetDuplicatedKeyError(xerror.NewError(int32(aiload.ErrCode_UserDuplicateKey)))
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
