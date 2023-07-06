package grpcapi

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/RomanIkonnikov93/niisva/internal/config"
	"github.com/RomanIkonnikov93/niisva/internal/models"
	pb "github.com/RomanIkonnikov93/niisva/internal/proto"
	ms "github.com/RomanIkonnikov93/niisva/internal/sync"
	"github.com/RomanIkonnikov93/niisva/pkg/logging"
	"github.com/fsnotify/fsnotify"

	"google.golang.org/protobuf/types/known/emptypb"
)

type WatcherServiceServer struct {
	pb.UnimplementedWatcherServer
	Users      models.Users
	UsersStore Repository
	cfg        *config.Config
	logger     *logging.Logger
}

type Repository interface {
	Add(ctx context.Context, path string) error
	GetAll(ctx context.Context) (map[string]struct{}, error)
}

func InitServices(ctx context.Context, cfg *config.Config, logger *logging.Logger, rep Repository) (*WatcherServiceServer, error) {

	w := WatcherServiceServer{
		Users:      models.Users{Paths: make(map[string]struct{})},
		UsersStore: rep,
		cfg:        cfg,
		logger:     logger,
	}

	m, err := w.UsersStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(m) != 0 {
		w.Users.Paths = m
	}

	return &w, nil
}

func (w *WatcherServiceServer) Run() error {

	path := filepath.Clean(w.cfg.FileStoragePath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {

					fileType := filepath.Ext(event.String())

					if fileType == w.cfg.FileType+`"` || fileType == w.cfg.FileType+`~"` {

						var wg sync.WaitGroup

						for val := range w.Users.Paths {
							wg.Add(1)
							go func(val string) {
								err = ms.FilesSync(path, val, w.cfg.FileType)
								if err != nil {
									w.logger.Println(err)
								}
								wg.Done()
							}(val)
						}
						wg.Wait()
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				w.logger.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(path)
	if err != nil {
		return err
	}

	<-make(chan struct{})

	return nil
}

func (w *WatcherServiceServer) AddUserRemoteDirectoryPath(ctx context.Context, in *pb.Record) (*emptypb.Empty, error) {

	val := filepath.Clean(in.Path)

	w.Users.Paths[val] = struct{}{}

	err := w.UsersStore.Add(ctx, val)
	if err != nil {
		return nil, err
	}

	out := &emptypb.Empty{}
	return out, nil
}
