package services

import (
	"fmt"

	"syncdata/pkg/repository"
	"syncdata/pkg/repository/entities"

	"github.com/ahmetb/go-linq/v3"
)

type GameRecordService struct {
	repo        *repository.GameRecordRepository
	allUserRepo *repository.AllUserRepository
}

func NewGameRecordService(repo *repository.GameRecordRepository, allUserRepo *repository.AllUserRepository) *GameRecordService {
	return &GameRecordService{
		repo:        repo,
		allUserRepo: allUserRepo,
	}
}

func (s *GameRecordService) Get(startTime string, endTime string) ([]entities.GameRecordEntity, error) {
	tables, err := s.repo.Exists(startTime)
	if err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, nil
	}
	sqlStr := fmt.Sprintf(`
		SET SESSION group_concat_max_len = 1073741824;	
		select
		CONCAT('SELECT ChannelId,Accounts,KindID,ServerID,CellScore,Profit,GameEndTime as CreateTime FROM ',
		TABLE_SCHEMA,'.', TABLE_NAME,' WHERE GameEndTime  BETWEEN ''%s'' AND ''%s''' )  INTO  @GameParameterNotMatchGameRecordDBName
		from information_schema.TABLES 
		where TABLE_SCHEMA ='dzpk_record' AND TABLE_NAME =concat('gameRecord', DATE_FORMAT( '%s' , '%%Y%%m%%d' ) );
		SELECT CONCAT(
			'SELECT records.*
					,CASE DATEDIFF(records.CreateTime ,accounts.createdate) WHEN 0 THEN 1 ELSE 0 END IsNew
					,accounts.lastlogintime as LastLoginTime
					,accounts.createdate as RegisterTime
					,0 as HistoryCellScore
					,0 as HistoryProfit
					,0 as HistoryGameNum
			FROM ('
			,  if(  @GameParameterNotMatchGameRecordDBName is null , '  ' , concat( @GameParameterNotMatchGameRecordDBName, ' \n UNION ALL   ' )   )  
			,GROUP_CONCAT(CONCAT('SELECT ChannelId,Accounts,KindID,ServerID,CellScore,Profit,GameEndTime  as CreateTime FROM ',it.TABLE_SCHEMA,'.',it.TABLE_NAME,' WHERE GameEndTime BETWEEN ''%s'' AND ''%s''') separator ' UNION ALL ')
			,') records
			INNER JOIN game_api.accounts accounts ON records.Accounts = accounts.account')
			INTO @v_sql
			FROM ( SELECT * FROM (SELECT t.* FROM KYDB_NEW.GameInfo as game INNER JOIN information_schema.TABLES as t ON t.TABLE_SCHEMA REGEXP CONCAT('^',game.GameParameter,'_record$')) as gameTemp 
				WHERE DATEDIFF(STR_TO_DATE('%s','%%Y-%%m-%%d'),STR_TO_DATE(gameTemp.TABLE_NAME,'gameRecord%%Y%%m%%d')) = 0) AS it;
			PREPARE stmt FROM @v_sql;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
		`, startTime, endTime, startTime, startTime, endTime, startTime)
	//取得所有GameRecord資料
	gameRecords, err := s.repo.QueryGameRecord(sqlStr)
	if err != nil {
		return nil, err
	}
	//取得聚合過後的AllUser資料
	allUser, allUserErr := s.allUserRepo.QueryAllUserRecord()

	if allUserErr != nil {
		return nil, allUserErr
	}
	//LeftJoin
	var leftJoinRecords []entities.GameRecordEntity
	linq.From(gameRecords).
		GroupJoinT(linq.From(allUser),
			func(gameRecord entities.GameRecordEntity) string { return gameRecord.Accounts },
			func(allUser entities.GroupByAllUserEntity) string { return allUser.Account },
			func(gameRecord entities.GameRecordEntity, allUser []entities.GroupByAllUserEntity) linq.KeyValue {
				return linq.KeyValue{Key: gameRecord, Value: allUser}
			},
		).SelectIndexedT(func(index int, obj linq.KeyValue) entities.GameRecordEntity {
		tmpGameRecord := obj.Key.(entities.GameRecordEntity)
		if allUserRecords, ok := obj.Value.([]entities.GroupByAllUserEntity); ok {
			if len(allUserRecords) > 0 {
				for _, allUser := range allUserRecords {
					tmpGameRecord.HistoryCellScore = allUser.CellScore
					tmpGameRecord.HistoryProfit = allUser.WinGold + allUser.LostGold
					tmpGameRecord.HistoryGameNum = allUser.WinNum + allUser.LostNum
				}
			}
		}
		return tmpGameRecord
	}).ToSlice(&leftJoinRecords)

	return leftJoinRecords, nil
}
