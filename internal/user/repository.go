package user

// Repository is the interface that provides user storage methods.
// It defines the contract that our business logic (usecase) will use.
type Repository interface {
	// We will define methods here in the next checkpoint, for example:
	// CreateUser(ctx context.Context, user *domain.User) error
	// GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
