package main

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"log"
	"time"
)

type Project struct {
	Id   string
	Name string
}
type ProjectRes struct {
	Links struct {
		Self     string      `json:"self"`
		Previous interface{} `json:"previous"`
		Next     interface{} `json:"next"`
	} `json:"links"`
	Projects []struct {
		IsDomain    bool   `json:"is_domain"`
		Description string `json:"description"`
		Links       struct {
			Self string `json:"self"`
		} `json:"links"`
		Tags     []interface{} `json:"tags"`
		Enabled  bool          `json:"enabled"`
		ID       string        `json:"id"`
		ParentID string        `json:"parent_id"`
		DomainID string        `json:"domain_id"`
		Name     string        `json:"name"`
	} `json:"projects"`
}
type TokenRes struct {
	Token struct {
		IssuedAt  time.Time `json:"issued_at"`
		AuditIds  []string  `json:"audit_ids"`
		Methods   []string  `json:"methods"`
		ExpiresAt time.Time `json:"expires_at"`
		User      struct {
			PasswordExpiresAt interface{} `json:"password_expires_at"`
			Domain            struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
	} `json:"token"`
}
type RoleRes struct {
	Token struct {
		IsDomain bool     `json:"is_domain"`
		Methods  []string `json:"methods"`
		Roles    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
		ExpiresAt time.Time `json:"expires_at"`
		Project   struct {
			Domain struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Catalog []struct {
			Endpoints []struct {
				RegionID  string `json:"region_id"`
				URL       string `json:"url"`
				Region    string `json:"region"`
				Interface string `json:"interface"`
				ID        string `json:"id"`
			} `json:"endpoints"`
			Type string `json:"type"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"catalog"`
		User struct {
			PasswordExpiresAt interface{} `json:"password_expires_at"`
			Domain            struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		AuditIds []string  `json:"audit_ids"`
		IssuedAt time.Time `json:"issued_at"`
	} `json:"token"`
}

// 得到临时token
func TokensTem(name, password string) (string, string) {
	request := gorequest.New()
	resp, body, errs := request.Post(KeystoneIp + "/v3/auth/tokens").
		Send(`{
    "auth": {
        "identity": {
            "methods": [
                "password"
            ],
            "password": {
                "user": {
                    "name":"` + name + `",
                    "domain": {
                        "name": "Default"
                    },
                    "password":"` + password + `"
                }
            }
        }
    }
}`).End()
	if errs != nil {
		log.Fatal(errs)
	}
	var tokenRes TokenRes
	err := json.Unmarshal([]byte(body), &tokenRes)
	if err != nil {
		log.Fatal("json error ", err)
	}
	return resp.Header.Get("X-Subject-Token"), tokenRes.Token.User.ID
}

// 获取用户的项目
func GetProjects(token string) []Project {
	request := gorequest.New()
	_, body, errs := request.Get(KeystoneIp+"/v3/auth/projects").
		Set("X-Auth-Token", token).End()
	if errs != nil {
		log.Fatal(errs)
	}
	//fmt.Println(body)
	var projectRes ProjectRes
	err := json.Unmarshal([]byte(body), &projectRes)
	if err != nil {
		log.Fatal("json error ", err)
	}
	projects := make([]Project, len(projectRes.Projects))
	for i := 0; i < len(projectRes.Projects); i++ {
		projects[i].Id = projectRes.Projects[i].ID
		projects[i].Name = projectRes.Projects[i].Name
	}
	return projects
}

// 获得某个用户在某个项目的角色
func GetRole(projectId string, token string) string {
	request := gorequest.New()
	_, body, errs := request.Post(KeystoneIp + "/v3/auth/tokens").
		Send(`{
    "auth": {
        "identity": {
            "methods": [
                "token"
            ],
            "token": {
                "id": "` + token + `"
            }
        },
        "scope": {
            "project": {
                "id": "` + projectId + `"
            }
        }
    }
}`).End()
	if errs != nil {
		log.Fatal(errs)
	}
	var roleRes RoleRes
	err := json.Unmarshal([]byte(body), &roleRes)
	if err != nil {
		log.Fatal("json error ", err)
	}
	return roleRes.Token.Roles[0].Name
}
