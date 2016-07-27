package session

var sessionMap map[string]string

func getSessionMap() map[string]string {
	if sessionMap == nil {
		sessionMap = make(map[string]string)
	}
	return sessionMap
}

func SetAttribute(key, value string) {
	getSessionMap()[key] = value
}

func GetAttribute(key string) string {
	return getSessionMap()[key]
}

func ContainsKey(key string) bool {
	_, logged := getSessionMap()[key]
	return logged
}

func RemoveAttribute(key string) {
	delete(sessionMap, key)
}
