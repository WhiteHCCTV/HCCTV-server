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
  switch ( os.Getenv("ENV") ){
    case "local":
      return "localhost:"+os.Getenv("LOCAL_PORT")
    case "dev":
      return "localhost:"+os.Getenv("DEV_PORT")
    case "prod":
      return os.Getenv("AWS")+":"+os.Getenv("AWS_PORT")
    }
    return "Invalid env"
  }
