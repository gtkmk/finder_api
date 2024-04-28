package follow

import (
	"github.com/gtkmk/finder_api/core/domain/followDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

const checkPointManageFollowUserTransactionNameConst = "manageFollowUser"

type CreateManageFollowUser struct {
	FollowDatabase repositories.FollowRepository
	Follow         *followDomain.Follow
	Following      bool
	Transaction    port.ConnectionInterface
	CustomError    port.CustomErrorInterface
}

func NewCreateManageFollowUser(
	followDatabase repositories.FollowRepository,
	follow *followDomain.Follow,
	transaction port.ConnectionInterface,
	customError port.CustomErrorInterface,
) *CreateManageFollowUser {
	return &CreateManageFollowUser{
		FollowDatabase: followDatabase,
		Follow:         follow,
		Transaction:    transaction,
		CustomError:    customError,
	}
}

func (createManageFollowUser *CreateManageFollowUser) Execute() (bool, error) {
	if err := createManageFollowUser.Transaction.SavePoint(checkPointManageFollowUserTransactionNameConst); err != nil {
		return false, createManageFollowUser.CustomError.ThrowError(err.Error())
	}

	err := createManageFollowUser.manageFollow()
	if err != nil {
		if rollbackErr := createManageFollowUser.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return false, rollbackErr
		}
		return false, err
	}

	if err := createManageFollowUser.Transaction.Commit(); err != nil {
		if rollbackErr := createManageFollowUser.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return false, rollbackErr
		}
		return false, createManageFollowUser.CustomError.ThrowError(err.Error())
	}

	return createManageFollowUser.Following, nil
}

func (createManageFollowUser *CreateManageFollowUser) manageFollow() error {
	alreadyFollowing, existingFollowInfo, err := createManageFollowUser.FollowDatabase.ConfirmExistingFollow(createManageFollowUser.Follow)
	if err != nil {
		return err
	}

	if alreadyFollowing {
		return createManageFollowUser.removeFollow(existingFollowInfo.Id)
	}

	return createManageFollowUser.createFollow()
}

func (createManageFollowUser *CreateManageFollowUser) removeFollow(existingFollowInfoId string) error {
	createManageFollowUser.Following = false
	return createManageFollowUser.FollowDatabase.RemoveFollow(existingFollowInfoId)
}

func (createManageFollowUser *CreateManageFollowUser) createFollow() error {
	createManageFollowUser.Following = true
	return createManageFollowUser.FollowDatabase.CreateFollow(createManageFollowUser.Follow)
}

func (createManageFollowUser *CreateManageFollowUser) rollbackToSavePointAndCommit() error {
	if transactErr := createManageFollowUser.Transaction.RollbackTo(checkPointManageFollowUserTransactionNameConst); transactErr != nil {
		return createManageFollowUser.CustomError.ThrowError(transactErr.Error())
	}

	if commitErr := createManageFollowUser.Transaction.Commit(); commitErr != nil {
		return createManageFollowUser.CustomError.ThrowError(commitErr.Error())
	}

	return nil
}
