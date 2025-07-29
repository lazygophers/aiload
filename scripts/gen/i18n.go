package main

import (
	"github.com/boltdb/bolt"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/i18n"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils"
	"github.com/lazygophers/utils/cryptox"
)

var destLangs = []*i18n.Language{
	i18n.MustParseLanguage("zh"),
	i18n.MustParseLanguage("zh-cn"),
	i18n.MustParseLanguage("zh-hans"),
	i18n.MustParseLanguage("en"),
}

func GenGoI18n() error {
	c := &i18n.TransacteConfig{
		SrcFile:            "./internal/state/localize/zh.yaml",
		SrcLang:            i18n.MustParseLanguage("zh"),
		Langs:              destLangs,
		OverwriteKeyPrefix: nil,
		Overwrite:          false,
		Localizer:          utils.MustOk(i18n.GetLocalizer(".yaml")),
		Translator:         utils.MustOk(i18n.GetTranslator(i18n.TransacteType(state.Config.I18n.Translator))),
	}

	dstLocalize, err := i18n.LoadLocalize(c.SrcFile, c.Localizer)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = codegen.GenerateI18nConst(dstLocalize, "internal/state/i18n_const.gen.go")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	cache, err := bolt.Open("scripts/gen/go.i18n.cache", 0600, &bolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		ReadOnly:        false,
		MmapFlags:       0,
		InitialMmapSize: 0,
	})
	if err != nil {
		return err
	}
	defer cache.Close()

	err = genI18n(c, cache)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GenJsI18n() error {
	c := &i18n.TransacteConfig{
		SrcFile:            "./assets/src/assets/localize/zh.yaml",
		SrcLang:            i18n.MustParseLanguage("zh"),
		Langs:              destLangs,
		OverwriteKeyPrefix: nil,
		Overwrite:          false,
		Localizer:          utils.MustOk(i18n.GetLocalizer(".yaml")),
		Translator:         utils.MustOk(i18n.GetTranslator(i18n.TransacteType(state.Config.I18n.Translator))),
	}

	cache, err := bolt.Open("scripts/gen/js.i18n.cache", 0600, &bolt.Options{
		Timeout:         0,
		NoGrowSync:      false,
		ReadOnly:        false,
		MmapFlags:       0,
		InitialMmapSize: 0,
	})
	if err != nil {
		return err
	}
	defer cache.Close()

	err = genI18n(c, cache)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func genI18n(c *i18n.TransacteConfig, cache *bolt.DB) error {
	tx, err := cache.Begin(true)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	i18n.RegisterBeforeTranslate(func(lang *i18n.Language, key, source, dest string) (skip bool) {
		if dest == "" {
			return false
		}

		bucket, err := tx.CreateBucketIfNotExists([]byte(lang.Lang))
		if err != nil {
			tx.Rollback()
			log.Panicf("err:%v", err)
			return
		}

		skip = string(bucket.Get([]byte(key))) == cryptox.Md5(source)

		return skip
	})

	i18n.RegisterAfterTranslate(func(lang *i18n.Language, key, source, dest string) {
		if dest == "" {
			return
		}

		bucket, err := tx.CreateBucketIfNotExists([]byte(lang.Lang))
		if err != nil {
			tx.Rollback()
			log.Panicf("err:%v", err)
			return
		}

		err = bucket.Put([]byte(key), []byte(cryptox.Md5(source)))
		if err != nil {
			tx.Rollback()
			log.Panicf("err:%v", err)
			return
		}
	})

	err = i18n.Translate(c)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = cache.Sync()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
