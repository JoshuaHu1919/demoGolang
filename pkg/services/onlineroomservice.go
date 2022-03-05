package services

import (
	"fmt"
	"strconv"
	"strings"

	"syncdata/internal/utils"
	"syncdata/pkg/repository"
	"syncdata/pkg/repository/entities"
	"syncdata/pkg/services/contract"

	"github.com/ahmetb/go-linq/v3"
	"github.com/rs/zerolog/log"
)

type OnlineRoomService struct {
	repo *repository.OnlineRoomRepository
}

func NewOnlineRoomService(repo *repository.OnlineRoomRepository) *OnlineRoomService {
	return &OnlineRoomService{
		repo: repo,
	}
}

func (s *OnlineRoomService) GetOnlineRoomRecord(date string) ([]contract.OnlineRoomContract, error) {
	rep, err := s.repo.QueryOnlineRoomRecord(date)
	if err != nil {
		return nil, err
	}
	//資料聚合
	var linqResults []linq.Group
	linq.From(rep).GroupByT(
		func(ele entities.OnlineRoomEntity) string { return ele.RoomId },
		func(ele entities.OnlineRoomEntity) entities.OnlineRoomEntity {

			return entities.OnlineRoomEntity{
				Id:         ele.Id,
				RoomId:     ele.RoomId,
				Value:      ele.Value,
				IP:         ele.IP,
				CreateTime: ele.CreateTime,
			}
		},
	).ToSlice(&linqResults)

	var resultRecords []contract.OnlineRoomContract
	for _, onlineRooms := range linqResults {
		var resultRecord contract.OnlineRoomContract
		paresInt, paresIntErr := strconv.Atoi(fmt.Sprintf("%v", onlineRooms.Key))
		if paresIntErr != nil {
			log.Err(paresIntErr).Msgf("Failed convert RoomId to Int")
		}
		resultRecord.RoomId = paresInt
		resultRecord.CreateDate = strings.Replace(date, "-", "", -1)
		for _, v := range onlineRooms.Group {
			newOnlineRoomContract, ok := v.(entities.OnlineRoomEntity)
			if !ok {
				return nil, fmt.Errorf("cast to entities.OnlineRoomEntity failed")
			}

			resultRecord.Value += float64(newOnlineRoomContract.Value)

		}
		resultRecords = append(resultRecords, resultRecord)
	}

	//聚合後的數值除於1440
	var calculateRecords []contract.OnlineRoomContract
	for _, resultRooms := range resultRecords {
		calculateRecords = append(calculateRecords, contract.OnlineRoomContract{
			RoomId:     resultRooms.RoomId,
			Value:      utils.Round(resultRooms.Value/1440, 4),
			CreateDate: resultRooms.CreateDate,
		})
	}
	return calculateRecords, nil
}
