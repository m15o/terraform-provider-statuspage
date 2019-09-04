package statuspage

import (
	"errors"
	"strings"
)

type pageResourceId struct {
	pageId     string
	resourceId string
}

func parsePageResourceId(id string) (*pageResourceId, error) {
	ids := strings.SplitN(id, "/", 2)
	if len(ids) != 2 {
		return nil, errors.New("id is not formatted properly; id should be '$page_id/$component_id', but: " + id)
	}
	return &pageResourceId{ids[0], ids[1]}, nil
}
