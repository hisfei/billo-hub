package api

import (
	"billohub/internal/skill"
	"billohub/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListAllSkillsResponse defines the response structure for listing all skills.
type ListAllSkillsResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

// ListAllSkills returns the global skill pool.
func (h *APIHandler) ListAllSkills(c *gin.Context) {
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	registeredSkills := skill.GetAllRegistered()
	skillList := make([]ListAllSkillsResponse, 0)

	for name, tempSkill := range registeredSkills {
		var t ListAllSkillsResponse
		t.Id = name
		t.Name = tempSkill.GetDescName()
		t.Desc = tempSkill.GetDescription()
		t.Status = "ENABLED"
		skillList = append(skillList, t)
	}
	res.Body = skillList
	c.JSON(http.StatusOK, res)
}
