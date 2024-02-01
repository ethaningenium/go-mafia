package services

func ExtractToken(authorizationHeader string) string {
	// Проверка, содержит ли заголовок префикс "Bearer "
	const bearerPrefix = "Bearer "
	if len(authorizationHeader) > len(bearerPrefix) && authorizationHeader[:len(bearerPrefix)] == bearerPrefix {
		return authorizationHeader[len(bearerPrefix):]
	}
	return authorizationHeader
}