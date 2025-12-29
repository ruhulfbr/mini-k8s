package repositories

import "github.com/ruhulfbr/mini-k8s/internal/datastore"

type Repositories struct {
	//ApplicationRepository user.Repository
	//ServiceRepository post.Repository
	//PodRepository post.Repository
}

func InitRepositories(ds *datastore.Datastore) *Repositories {
	return &Repositories{
		//UserRepository: userRepo.NewUserRepository(db),
		//PostRepository: postRepo.NewPostRepository(db),
	}
}
