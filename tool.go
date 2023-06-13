package jsonschemaline

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/pkg/errors"
)

func JsonMerge(first string, second string, more ...string) (merge string, err error) {
	combinedPatch := []byte(second)
	for _, patch := range more {
		combinedPatch, err = jsonpatch.MergeMergePatches(combinedPatch, []byte(patch))
		if err != nil {
			err = errors.WithStack(err)
			return "", err
		}

	}
	mb, err := jsonpatch.MergePatch([]byte(first), combinedPatch)
	if err != nil {
		err = errors.WithStack(err)
		return "", err
	}
	merge = string(mb)
	return merge, err
}
