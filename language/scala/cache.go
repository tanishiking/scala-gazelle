package scala

import (
	"log"
	"time"

	"github.com/stackb/scala-gazelle/pkg/protobuf"
)

func (sl *scalaLang) readCacheFile() error {
	t1 := time.Now()

	if err := protobuf.ReadFile(sl.cacheFileFlagValue, sl.cache); err != nil {
		return err
	}
	for _, rule := range sl.cache.Rules {
		if err := sl.sourceProvider.ProvideRule(rule); err != nil {
			return err
		}
	}

	t2 := time.Since(t1).Round(1 * time.Millisecond)

	if debug {
		log.Printf("Read cache %s (%d rules) %v", sl.cacheFileFlagValue, len(sl.cache.Rules), t2)
	}
	return nil
}

func (sl *scalaLang) writeCacheFile() error {
	sl.cache.PackageCount = int32(len(sl.packages))
	sl.cache.Rules = sl.sourceProvider.ProvidedRules()

	if debug {
		log.Printf("Wrote cache %s (%d rules)", sl.cacheFileFlagValue, len(sl.cache.Rules))
	}
	return protobuf.WriteFile(sl.cacheFileFlagValue, sl.cache)
}