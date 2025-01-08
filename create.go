package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a custom directory name for 'history'.")
		fmt.Println("Usage: go run main.go <custom-directory-name>")
		return
	}

	customDirName := strings.ToLower(os.Args[1])
	projectName := "washit-api"

	dirs := []string{
		"dto/model",
		"dto/request",
		"dto/resource",
		"handler",
		"repository",
		"routes",
		"service",
	}

	for _, dir := range dirs {
		path := fmt.Sprintf("./internal/%s/%s", customDirName, dir)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", path, err)
			return
		}
	}

	createGoFile := func(filePath, packageName, content string) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("package %s\n\n%s", packageName, content))
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", filePath, err)
		}
	}

	createGoFile(fmt.Sprintf("./internal/%s/dto/model/%s-model.go", customDirName, customDirName), customDirName+"Model", `
type `+strings.Title(customDirName)+`Model struct {
 ID   int    `+"`json:\"id\"`"+`
 Name string `+"`json:\"name\"`"+`
}`)

	createGoFile(fmt.Sprintf("./internal/%s/dto/request/%s-request.go", customDirName, customDirName), customDirName+"Request", `
type `+strings.Title(customDirName)+`Request struct {
 ID   int    `+"`json:\"id\"`"+`
 Name string `+"`json:\"name\"`"+`
}`)

	createGoFile(fmt.Sprintf("./internal/%s/dto/resource/%s-resource.go", customDirName, customDirName), customDirName+"Resource", `
type `+strings.Title(customDirName)+`Resource struct {
 ID   int    `+"`json:\"id\"`"+`
 Name string `+"`json:\"name\"`"+`
}`)

	createGoFile(fmt.Sprintf("./internal/%s/handler/%s.go", customDirName, customDirName), customDirName, `
import (
 "github.com/gin-gonic/gin"
 "net/http"

 `+customDirName+`Service "`+projectName+`/internal/`+customDirName+`/service"
 "`+projectName+`/pkg/redis"
 "`+projectName+`/pkg/response"
)

type `+strings.Title(customDirName)+`Handler struct {
 service `+customDirName+`Service.I`+strings.Title(customDirName)+`Service
 cache   redis.IRedis
}

func New`+strings.Title(customDirName)+`Handler(service `+customDirName+`Service.I`+strings.Title(customDirName)+`Service, cache redis.IRedis) *`+strings.Title(customDirName)+`Handler {
 return &`+strings.Title(customDirName)+`Handler{
  service: service,
  cache:   cache,
 }
}

func (h *`+strings.Title(customDirName)+`Handler) New(c *gin.Context) {
 response.Success(c, http.StatusOK, "Successfully make", nil, nil)
}`)

	createGoFile(fmt.Sprintf("./internal/%s/repository/%s-repository.go", customDirName, customDirName), customDirName+"Repository", `
import (
 "`+projectName+`/pkg/db/dbs"
)

type I`+strings.Title(customDirName)+`Repository interface {
}

type `+strings.Title(customDirName)+`Repository struct {
 db dbs.IDatabase
}

func New`+strings.Title(customDirName)+`Repository(db dbs.IDatabase) *`+strings.Title(customDirName)+`Repository {
 return &`+strings.Title(customDirName)+`Repository{db: db}
}`)

	createGoFile(fmt.Sprintf("./internal/%s/routes/routes.go", customDirName), customDirName+"Routes", `
import (
 "github.com/gin-gonic/gin"
 "github.com/go-playground/validator"

 `+customDirName+` "`+projectName+`/internal/`+customDirName+`/handler"
 `+customDirName+`Repository "`+projectName+`/internal/`+customDirName+`/repository"
 `+customDirName+`Service "`+projectName+`/internal/`+customDirName+`/service"
 "`+projectName+`/pkg/db/dbs"
 "`+projectName+`/pkg/middleware"
 "`+projectName+`/pkg/redis"
)

func Main(r *gin.RouterGroup, db dbs.IDatabase, cache redis.IRedis, validator *validator.Validate) {
 repository := `+customDirName+`Repository.New`+strings.Title(customDirName)+`Repository(db)
 service := `+customDirName+`Service.New`+strings.Title(customDirName)+`Service(repository, validator)
 handler := `+customDirName+`.New`+strings.Title(customDirName)+`Handler(service, cache)

 authMiddleware := middleware.JWTAuth()
 adminAuthMiddleware := middleware.JTWAuthAdmin()

 r.GET("/new", adminAuthMiddleware, handler.New)
 r.GET("/new", authMiddleware, handler.New)
}`)

	createGoFile(fmt.Sprintf("./internal/%s/service/%s-service.go", customDirName, customDirName), "service", `
import (
 "github.com/go-playground/validator"

 `+customDirName+`Repository "`+projectName+`/internal/`+customDirName+`/repository"
)

type I`+strings.Title(customDirName)+`Service interface {
}

type `+strings.Title(customDirName)+`Service struct {
 repository `+customDirName+`Repository.I`+strings.Title(customDirName)+`Repository
 validator  *validator.Validate
}

func New`+strings.Title(customDirName)+`Service(
 repository `+customDirName+`Repository.I`+strings.Title(customDirName)+`Repository, validator *validator.Validate) *`+strings.Title(customDirName)+`Service {
 return &`+strings.Title(customDirName)+`Service{
  repository: repository,
  validator:  validator,
 }
}`)

	fmt.Printf("'%s' have been created inside './internal'.\n", customDirName)
}
