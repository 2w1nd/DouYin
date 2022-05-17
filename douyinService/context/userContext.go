package context

type UserContext struct {
	id uint64
}

func token2UserContext(token string) UserContext {
	return UserContext{}
}
