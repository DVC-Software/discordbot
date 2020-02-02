package execution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DVC-Software/discordbot/config"
)

type MemberRequest struct {
	Name          string
	SlackUserID   string
	DiscordUserID string
	CreatedFrom   string
}

type MemberInfo struct {
	ID             uint
	Name           string
	Positions      []string
	TrainingStatus int
	IsStaff        bool
}

var endpoint = config.DVCApiServerEndpoint

func CreateMember(args []string) (string, string) {
	requestBody := MemberRequest{Name: args[0], SlackUserID: "", DiscordUserID: args[1], CreatedFrom: "discord"}
	body, _ := json.Marshal(requestBody)
	resp, err := http.Post(endpoint+"/member/create", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "Failed to create member " + args[0], err.Error()
	}
	if resp.StatusCode != 200 {
		if resp.Body != nil {
			data, _ := ioutil.ReadAll(resp.Body)
			var respError Error
			json.Unmarshal(data, &respError)
			return "Failed to create member " + args[0], respError.ErrorMessage
		}
		return "Failed to create member " + args[0], "Unknown error, status: " + resp.Status
	}
	// parse returned obj
	data, _ := ioutil.ReadAll(resp.Body)
	var memberInfo MemberInfo
	json.Unmarshal(data, &memberInfo)
	if memberInfo.Name != args[0] {
		// debug
		return "Failed to create member " + args[0], "Invalid response data"
	}

	return resultToString(memberInfo.ID, memberInfo.Name, memberInfo.Positions[0]), ""
}

func resultToString(ID uint, name string, position string) string {
	return fmt.Sprintf("Member ID : %d\nName : %v\nPosition: %s\nFrom now on, I'll call you %v!", ID, name, position, name)
}

func IdentifyMember(id string) (bool, MemberInfo) {
	var memberInfo MemberInfo
	resp, err := http.Get(endpoint + "/member/info/" + id)
	if err != nil || resp.StatusCode != 200 {
		if resp.Body != nil {
			data, _ := ioutil.ReadAll(resp.Body)
			var respError Error
			json.Unmarshal(data, &respError)
			fmt.Println(respError.ErrorMessage)
		}
		return false, memberInfo
	}
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	json.Unmarshal(data, &memberInfo)
	if memberInfo.ID == 0 || memberInfo.Name == "" {
		fmt.Println(memberInfo.ID, memberInfo.Name)
		return false, memberInfo
	}
	return true, memberInfo
}
