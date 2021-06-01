package eth

import (
	"github.com/pkg/errors"
)

// translateTopics takes in an ethReq's topics field and translates it to a it's equivalent QtumReq
// topics (optional) has a max lenght of 4
func TranslateTopics(ethTopics []interface{}) ([]interface{}, error) {

	var topics []interface{}

	if len(ethTopics) > 4 {
		return nil, errors.Errorf("invalid number of topics. Logs have a max of 4 topics.")
	}

	for _, topic := range ethTopics {
		switch topic.(type) {
		case []interface{}:
			topic, err := TranslateTopics(topic.([]interface{}))
			if err != nil {
				return nil, err
			}
			topics = append(topics, topic)
		case string:
			// topics = append(topics, utils.RemoveHexPrefix(topic.(string)))
			topics = append(topics, topic.(string))
		case nil:
			topics = append(topics, nil)
		}
	}

	return topics, nil

}
