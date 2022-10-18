package valorant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
)

type (
	BaseResp struct {
		Status int     `json:"status"`
		Errors []Error `json:"errors"`
	}

	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Details string `json:"details"`
	}

	MatchesResp struct {
		BaseResp
		Data []MatchData `json:"data"`
	}

	AccountResp struct {
		BaseResp
		Data AccountDetails `json:"data"`
	}

	AccountDetails struct {
		PUuid  string `json:"puuid"`
		Region string `json:"region"`
		Name   string `json:"name"`
		Tag    string `json:"tag"`
	}

	MatchData struct {
		MetaData MatchMetaData `json:"metadata"`
	}

	MatchMetaData struct {
		Map              string `json:"map"`
		GameVersion      string `json:"game_version"`
		GameLength       int64  `json:"game_length"`
		GameStart        int64  `json:"game_start"`
		GameStartPatched string `json:"game_start_patched"`
		RoundsPlayed     int    `json:"rounds_played"`
		Mode             string `json:"mode"`
		Queue            string `json:"queue"`
		SeasonId         string `json:"season_id"`
		Platform         string `json:"platform"`
		MatchId          string `json:"matchid"`
		Region           string `json:"region"`
		Cluster          string `json:"cluster"`
	}

	MatchMetaDataPublish struct {
		Data        *MatchMetaData `json:"data"`
		DisplayName string         `json:"display_name"`
	}
)

const BaseUrl = "https://api.henrikdev.xyz/valorant"

func doRequest(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Authorization", valorantYaml.HenrikApiKey)
	return http.DefaultClient.Do(req)
}

func getValorantMatches(uuid string) ([]MatchData, error) {
	path := fmt.Sprintf("/v3/by-puuid/matches/%s/%s", valorantYaml.Region, uuid)
	resp, err := doRequest(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var matchesResp MatchesResp

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &matchesResp)
	if err != nil {
		logger.Errorf("error while parsing valorant matches response: %s", err)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	} else if len(matchesResp.Errors) > 0 {
		var apiErrors = make([]string, len(matchesResp.Errors))
		for i, apiError := range matchesResp.Errors {
			apiErrors[i] = fmt.Sprintf("(%d)%s", apiError.Code, apiError.Message)
		}
		return nil, fmt.Errorf("api errors: %s", strings.Join(apiErrors, ", "))
	} else if resp.StatusCode != 200 {
		return nil, &request.HttpError{
			Code:     resp.StatusCode,
			Status:   resp.Status,
			Response: resp,
		}
	} else {
		return matchesResp.Data, nil
	}
}

func getDisplayName(uuid string) (*AccountDetails, error) {
	path := fmt.Sprintf("/v1/by-puuid/account/%s", uuid)
	resp, err := doRequest(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var accountResp AccountResp
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &accountResp)
	if err != nil {
		logger.Errorf("error while parsing valorant matches response: %s", err)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	} else if len(accountResp.Errors) > 0 {
		var apiErrors = make([]string, len(accountResp.Errors))
		for i, apiError := range accountResp.Errors {
			apiErrors[i] = fmt.Sprintf("(%d)%s", apiError.Code, apiError.Message)
		}
		return nil, fmt.Errorf("api errors: %s", strings.Join(apiErrors, ", "))
	} else if resp.StatusCode != 200 {
		return nil, &request.HttpError{
			Code:     resp.StatusCode,
			Status:   resp.Status,
			Response: resp,
		}
	} else {
		return &accountResp.Data, nil
	}
}
