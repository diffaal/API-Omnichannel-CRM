package service

import (
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/presentation"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AgentService struct {
	userRepo        repository.IUserRepository
	interactionRepo repository.IinteractionRepository
}

type IAgentService interface {
	GetAgentList(map[string]interface{}) (map[string]interface{}, error)
}

func NewAgentService(userRepo repository.IUserRepository, interactionRepo repository.IinteractionRepository) *AgentService {
	agentService := AgentService{
		userRepo:        userRepo,
		interactionRepo: interactionRepo,
	}
	return &agentService
}

func (as *AgentService) GetAgentList(filters map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var agentDashboardList []presentation.AgentDashboardData

	agentList, err := as.userRepo.GetAgentList(filters)
	if (agentList == nil && err == nil) || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, enum.ERROR_DATA_NOT_FOUND

	} else if err != nil {
		return nil, err
	}

	for _, v := range agentList {
		agentData := presentation.AgentDashboardData{
			AgentId:   v.ID,
			AgentName: fmt.Sprintf("%s %s", v.FirstName, v.LastName),
		}

		activeInteractionCount, err := as.interactionRepo.GetActiveInteractionCount(v.ID)
		if err != nil {
			return nil, err
		}
		agentData.ActiveInteraction = activeInteractionCount

		interactionTodayCount, err := as.interactionRepo.GetInteractionHandledTodayCount(v.ID)
		if err != nil {
			return nil, err
		}
		agentData.InteractionHandled = interactionTodayCount

		agentDashboardList = append(agentDashboardList, agentData)
	}

	result["agents"] = agentDashboardList

	return result, nil
}
