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
	var member map[string]interface{}
	json.Unmarshal(data, &member)
	profile := member["Profile"].(map[string]interface{})
	discordUserID := member["DiscordUserID"].(map[string]interface{})
	if discordUserID["String"] != args[1] || profile["Name"] != args[0] || discordUserID["Valid"] != true {
		// debug
		fmt.Println(member["DiscordUserID"], profile["Name"])
		return "Failed to create member " + args[0], "Invalid response data"
	}

	return resultToString(member["ID"], profile), ""
}

func resultToString(ID interface{}, profile map[string]interface{}) string {
	return fmt.Sprintf("Member ID : %v\nName : %v\nFrom now on, I'll call you %v!", ID, profile["Name"], profile["Name"])
}
