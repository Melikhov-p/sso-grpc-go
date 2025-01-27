package file

import (
	"context"
	"fmt"

	"github.com/Melikhov-p/sso-grpc-go/internal/models"
	"github.com/Melikhov-p/sso-grpc-go/internal/storage"
	"go.uber.org/zap"
)

type FileStorage struct {
	log       *zap.Logger
	filePath  string
	users     map[int64]*models.User
	apps      map[int32]*models.App
	lastIndex int64
}

func NewFileStorage(log *zap.Logger, filePath string) *FileStorage {
	return &FileStorage{
		log:       log,
		filePath:  filePath,
		users:     make(map[int64]*models.User),
		apps:      make(map[int32]*models.App),
		lastIndex: 0,
	}
}

func (fs *FileStorage) CreateUser(_ context.Context, email string, passHash []byte) (int64, error) {
	fs.lastIndex++
	fs.users[fs.lastIndex] = &models.User{
		ID:       fs.lastIndex,
		Email:    email,
		PassHash: passHash,
	}
	return fs.lastIndex, nil
}

func (fs *FileStorage) GetUserByEmail(_ context.Context, email string) (*models.User, error) {
	for id, user := range fs.users {
		if user.Email == email {
			return fs.users[id], nil
		}
	}

	return nil, storage.ErrNotFound
}

// GetAppByID find app with provided ID in storage map "apps".
/* If app with ID does not exist - create new App with ID. */
func (fs *FileStorage) GetAppByID(_ context.Context, appID int32) (*models.App, error) {
	op := "storage.File.GetAppByID"

	var err error
	app := fs.apps[appID]

	if app == nil {
		fs.log.Debug("app with provided id does not exists", zap.Int("appID", int(appID)))

		app, err = fs.CreateApp(context.Background(), appID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return app, nil
}

func (fs *FileStorage) CreateApp(_ context.Context, appID int32) (*models.App, error) {
	newApp := &models.App{
		ID:        appID,
		Name:      fmt.Sprintf("App_%d", appID),
		SecretKey: "$2a$12$Bues8rdmFfS1QVc0XZI88eqlzFlAxQM.GWjUZhIfAQrnbaXbRKLVa",
	}

	fs.apps[appID] = newApp

	return newApp, nil
}
