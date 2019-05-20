package opsgenie

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func flattenOpsGenieTeamMembers(input []team.Member) []interface{} {
	members := make([]interface{}, 0, len(input))
	for _, inputMember := range input {
		outputMember := make(map[string]interface{})
		outputMember["username"] = inputMember.User
		outputMember["role"] = inputMember.Role

		members = append(members, outputMember)
	}

	return members
}
