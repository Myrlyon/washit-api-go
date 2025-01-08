#!/bin/bash

# Check if the ${CUSTOM_DIR_NAME} provided a custom directory name for "history"
if [ -z "$1" ]; then
  echo "Please provide a custom directory name for 'history'."
  echo "Usage: $0 <custom-directory-name>"
  exit 1
fi

# Convert the custom directory name to lowercase
CUSTOM_DIR_NAME=$(echo "$1" | tr '[:upper:]' '[:lower:]')

PROJECT_NAME="washit-api"

# Define the directory structure
DIRS=(
  "dto/model"
  "dto/request"
  "dto/resource"
  "handler"
  "repository"
  "routes"
  "service"
)

# Create the directory structure
for DIR in "${DIRS[@]}"; do
  mkdir -p "./internal/$CUSTOM_DIR_NAME/$DIR"
done

# Function to create Go files with content
create_go_file() {
  local FILE_PATH=$1
  local PACKAGE_NAME=$2
  local CONTENT=$3

  echo "package $PACKAGE_NAME" > "$FILE_PATH"
  echo "$CONTENT" >> "$FILE_PATH"
}

# Create basic Go files inside the directories with lowercase names and write content into them
create_go_file "./internal/$CUSTOM_DIR_NAME/dto/model/${CUSTOM_DIR_NAME}-model.go" "${CUSTOM_DIR_NAME}Model" \
"
type ${CUSTOM_DIR_NAME^}Model struct {
  ID   int    \`json:\"id\"\`
  Name string \`json:\"name\"\`
}"

#################################################################################################

create_go_file "./internal/$CUSTOM_DIR_NAME/dto/request/${CUSTOM_DIR_NAME}-request.go" "${CUSTOM_DIR_NAME}Request" \
"
type ${CUSTOM_DIR_NAME^}Request struct {
  ID   int    \`json:\"id\"\`
  Name string \`json:\"name\"\`
}"

#################################################################################################

create_go_file "./internal/$CUSTOM_DIR_NAME/dto/resource/${CUSTOM_DIR_NAME}-resource.go" "${CUSTOM_DIR_NAME}Resource" \
"
type ${CUSTOM_DIR_NAME^}Resource struct {
  ID   int    \`json:\"id\"\`
  Name string \`json:\"name\"\`
}"

#################################################################################################

create_go_file "./internal/$CUSTOM_DIR_NAME/handler/${CUSTOM_DIR_NAME}.go" "$CUSTOM_DIR_NAME" \
"import (
	\"github.com/gin-gonic/gin\"
	\"net/http\"

	${CUSTOM_DIR_NAME}Service \"${PROJECT_NAME}/internal/${CUSTOM_DIR_NAME}/service\"
	\"${PROJECT_NAME}/pkg/redis\"
	\"${PROJECT_NAME}/pkg/response\"
)

type ${CUSTOM_DIR_NAME^}Handler struct {
	service ${CUSTOM_DIR_NAME}Service.I${CUSTOM_DIR_NAME^}Service
	cache   redis.IRedis
}

func New${CUSTOM_DIR_NAME^}Handler(service ${CUSTOM_DIR_NAME}Service.I${CUSTOM_DIR_NAME^}Service, cache redis.IRedis) *${CUSTOM_DIR_NAME^}Handler {
	return &${CUSTOM_DIR_NAME^}Handler{
		service: service,
		cache:   cache,
	}
}

func (h *${CUSTOM_DIR_NAME^}Handler) New(c *gin.Context) {
	response.Success(c, http.StatusOK, \"Successfully make\", nil, nil)
}"

#################################################################################################

create_go_file "./internal/$CUSTOM_DIR_NAME/repository/${CUSTOM_DIR_NAME}-repository.go" "${CUSTOM_DIR_NAME}Repository" \
"import (
	\"${PROJECT_NAME}/pkg/db/dbs\"
)

type I${CUSTOM_DIR_NAME^}Repository interface {
}

type ${CUSTOM_DIR_NAME^}Repository struct {
	db dbs.IDatabase
}

func New${CUSTOM_DIR_NAME^}Repository(db dbs.IDatabase) *${CUSTOM_DIR_NAME^}Repository {
	return &${CUSTOM_DIR_NAME^}Repository{db: db}
}"


#################################################################################################

create_go_file "./internal/$CUSTOM_DIR_NAME/routes/routes.go" "${CUSTOM_DIR_NAME}Routes" \
"
import (
	\"github.com/gin-gonic/gin\"
	\"github.com/go-playground/validator\"

	${CUSTOM_DIR_NAME} \"${PROJECT_NAME}/internal/${CUSTOM_DIR_NAME}/handler\"
	${CUSTOM_DIR_NAME}Repository \"${PROJECT_NAME}/internal/${CUSTOM_DIR_NAME}/repository\"
	${CUSTOM_DIR_NAME}Service \"${PROJECT_NAME}/internal/${CUSTOM_DIR_NAME}/service\"
	\"${PROJECT_NAME}/pkg/db/dbs\"
	\"${PROJECT_NAME}/pkg/middleware\"
	\"${PROJECT_NAME}/pkg/redis\"
)

func Main(r *gin.RouterGroup, db dbs.IDatabase, cache redis.IRedis, validator *validator.Validate) {
	repository := ${CUSTOM_DIR_NAME}Repository.New${CUSTOM_DIR_NAME^}Repository(db)
	service := ${CUSTOM_DIR_NAME}Service.New${CUSTOM_DIR_NAME^}Service(repository, validator)
	handler := ${CUSTOM_DIR_NAME}.New${CUSTOM_DIR_NAME^}Handler(service, cache)

	authMiddleware := middleware.JWTAuth()
	adminAuthMiddleware := middleware.JTWAuthAdmin()

	r.GET(\"/new\", adminAuthMiddleware, handler.New)
	r.GET(\"/new\", authMiddleware, handler.New)
}"

#################################################################################################

create_go_file "./internal/$CUSTOM_DIR_NAME/service/${CUSTOM_DIR_NAME}-service.go" "service" \
"import (
	\"github.com/go-playground/validator\"

	${CUSTOM_DIR_NAME}Repository \"washit-api/internal/${CUSTOM_DIR_NAME}/repository\"
)

type I${CUSTOM_DIR_NAME^}Service interface {
}

type ${CUSTOM_DIR_NAME^}Service struct {
	repository ${CUSTOM_DIR_NAME}Repository.I${CUSTOM_DIR_NAME^}Repository
	validator  *validator.Validate
}

func New${CUSTOM_DIR_NAME^}Service(
	repository ${CUSTOM_DIR_NAME}Repository.I${CUSTOM_DIR_NAME^}Repository, validator *validator.Validate) *${CUSTOM_DIR_NAME^}Service {
	return &${CUSTOM_DIR_NAME^}Service{
		repository: repository,
		validator:  validator,
	}
}"

# Output success message
echo "'$CUSTOM_DIR_NAME' have been created inside './internal'."
