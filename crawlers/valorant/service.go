package valorant

import (
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"strings"
)

type (
	MatchesResp struct {
		Status int         `json:"status"`
		Data   []MatchData `json:"data"`
		Errors []Error     `json:"errors"`
	}

	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Details string `json:"details"`
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

func getValorantMatches(name, tag string) ([]MatchData, error) {
	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v3/matches/%s/%s/%s", valorantYaml.Region, name, tag)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Authorization", valorantYaml.HenrikApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var matchesResp MatchesResp
	err = json.NewDecoder(resp.Body).Decode(&matchesResp)
	if err != nil {
		return nil, err
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
