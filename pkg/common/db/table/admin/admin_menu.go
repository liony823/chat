package admin

import "context"

type AdminUserMenu struct {
	UserID string   `bson:"user_id"`
	Menus  []string `bson:"menus"`
}

type AdminMenu struct {
	Key          string `bson:"key"`
	Name         string `bson:"name"`
	Path         string `bson:"path"`
	Component    string `bson:"component"`
	Icon         string `bson:"icon"`
	Sort         int32  `bson:"sort"`
	Parent       string `bson:"parent"`
	Layout       bool   `bson:"layout"`
	HiddenInMenu bool   `bson:"hiddenInMenu"`
	Redirect     string `bson:"redirect"`
}

func (AdminMenu) TableName() string {
	return "admin_menu"
}

type AdminMenuInterface interface {
	Create(ctx context.Context, menus []*AdminMenu) error
	Update(ctx context.Context, key string, data map[string]any) error
	Delete(ctx context.Context, keys []string) error
	Take(ctx context.Context, key string) (*AdminMenu, error)
	List(ctx context.Context, parent string) ([]*AdminMenu, error)
	ListByKeys(ctx context.Context, keys []string) ([]*AdminMenu, error)
}

type AdminUserMenuInterface interface {
	Create(ctx context.Context, userMenus []*AdminUserMenu) error
	Update(ctx context.Context, userID string, menus []string) error
	Delete(ctx context.Context, userIDs []string) error
	Take(ctx context.Context, userID string) (*AdminUserMenu, error)
	List(ctx context.Context) ([]*AdminUserMenu, error)
}
