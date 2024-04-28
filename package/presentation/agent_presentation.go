package presentation

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AgentDashboardData struct {
	AgentId            string `json:"agent_id"`
	AgentName          string `json:"agent_name"`
	ActiveInteraction  int64  `json:"active_interaction"`
	InteractionHandled int64  `json:"interaction_handled"`
	Status             string `json:"status"`
}

func ParseGetListAgentFilters(c *gin.Context) (map[string]interface{}, error) {
	filters := make(map[string]interface{})

	agentIdsQuery := c.Query("agent_ids")
	keywordQuery := c.Query("keyword")
	pageQuery := c.Query("page")
	pageSizeQuery := c.Query("pageSize")

	if agentIdsQuery != "" {
		agentIds := strings.Split(agentIdsQuery, ",")
		filters["agent_ids"] = agentIds
	}

	if keywordQuery != "" {
		filters["keyword"] = keywordQuery
	}

	if pageQuery != "" {
		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			return nil, err
		}
		filters["page"] = page
	}

	if pageSizeQuery != "" {
		pageSize, err := strconv.Atoi(pageSizeQuery)
		if err != nil {
			return nil, err
		}
		filters["pageSize"] = pageSize
	}

	return filters, nil
}
