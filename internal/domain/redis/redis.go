package redis_domain

type RedisConnection struct {
	Host     string
	Port     int
	Password *string
}
