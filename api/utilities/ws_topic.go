package utilities

import (
	"errors"
	"likeminds-pandemonium/common"
	"strings"
)

// GetTopicSplit will decode topic and return split
func GetTopicSplit(topic string) ([]string, error) {
	if topic == "" || topic == "null" {
		return nil, errors.New(common.ErrorTopicMissing)
	}
	topicSplit := strings.Split(topic, ":")
	if len(topicSplit) <= 1 {
		return nil, errors.New(common.ErrorTopicInvalid)
	}
	return topicSplit, nil
}
