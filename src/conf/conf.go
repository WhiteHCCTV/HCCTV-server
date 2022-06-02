package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load("../.env")
    
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}
func GetEnv(key string) string{
	return os.Getenv(key)
}
func GetAddr() string{
    if (os.Getenv("ENV") == "dev"){
        return os.Getenv("HOST")+":"+os.Getenv("DOMAIN_PORT")
    } else {
        return os.Getenv("AWS")+":"+os.Getenv("AWS_PORT")
    }
    
}