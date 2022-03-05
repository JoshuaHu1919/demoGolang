package services

import (
	"time"

	configmanager "syncdata/internal/config"
	"syncdata/pkg/repository"
	"syncdata/pkg/services/contract"
)

type AccountsService struct {
	repo *repository.AccountsRepository
}

func NewAccountsService(repo *repository.AccountsRepository) *AccountsService {
	return &AccountsService{
		repo: repo,
	}
}

func (s *AccountsService) GetAccountsRecord() ([]contract.AccountsContract, error) {
	rep, err := s.repo.QueryAccountsRecord()
	if err != nil {
		return nil, err
	}
	utc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil, err
	}
	var resultRecords []contract.AccountsContract
	for _, player := range rep {
		calculateRecords := &contract.AccountsContract{
			Platform:      configmanager.GlobalConfig.SyncDataPlatform.Platform,
			Account:       player.Account,
			Agent:         player.Agent,
			RegisterTime:  player.CreateDate,
			LastLoginTime: player.LastLoginTime,
			StatisticDate: time.Now().UTC().In(utc).AddDate(0, 0, -1).Format("20060102"),
			ExpireTime:    time.Now().UTC().In(utc).Add(24 * time.Hour * 39),
		}
		resultRecords = append(resultRecords, *calculateRecords)
	}

	return resultRecords, nil
}
