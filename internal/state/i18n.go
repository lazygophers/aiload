package state

import (
	"embed"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/routine"
	"github.com/lazygophers/utils/runtime"
	"io/fs"
	"os"
	"path/filepath"
)

type LocalizeFs struct {
}

func (p *LocalizeFs) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(runtime.ExecDir(), name))
}

func (p *LocalizeFs) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(runtime.ExecDir(), name))
}

//go:embed localize/*.yaml
var i18nFs embed.FS

func LoadI18n() (err error) {
	State.I18n = i18n.DefaultI18n

	err = State.I18n.LoadLocalizesWithFs("localize", i18nFs)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	//State.I18n.SetDefaultLang(env.Language())
	State.I18n.SetDefaultLang("zh")

	xerror.SetI18n(i18n.NewI18nForXerror(State.I18n))

	routine.AddBeforeRoutine(func(baseGid, currentGid int64) {
		i18n.SetLanguage(i18n.GetLanguage(baseGid), currentGid)
	})
	routine.AddAfterRoutine(func(currentGid int64) {
		i18n.DelLanguage(currentGid)
	})

	return nil
}

func Localize(key string, args ...interface{}) string {
	return State.I18n.Localize(key, args...)
}

func LocalizeWithLang(lang string, key string, args ...interface{}) string {
	return State.I18n.LocalizeWithLang(lang, key, args...)
}
