package context

type UserContext struct {
	Id uint64
}

func token2UserContext(token string) UserContext {
	return UserContext{}
}
