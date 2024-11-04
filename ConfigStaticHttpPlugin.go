package configstatic_plugin

import (
        "github.com/LeakIX/l9format"
        "regexp"
        "strings"
)

type ConfigStaticHttpPlugin struct {
        l9format.ServicePluginBase
}

func (ConfigStaticHttpPlugin) GetVersion() (int, int, int) {
        return 0, 0, 1
}

func (ConfigStaticHttpPlugin) GetRequests() []l9format.WebPluginRequest {
        return []l9format.WebPluginRequest{{
                Method:  "GET",
                Path:    "/static../.git/config",
                Headers: map[string]string{},
                Body:    []byte(""),
        }}
}

func (ConfigStaticHttpPlugin) GetName() string {
        return "ConfigStaticHttpPlugin"
}

func (ConfigStaticHttpPlugin) GetStage() string {
        return "open"
}

func (plugin ConfigStaticHttpPlugin) Verify(request l9format.WebPluginRequest, response l9format.WebPluginResponse, event *l9format.L9Event, options map[string]string) (hasLeak bool) {
        if !request.EqualAny(plugin.GetRequests()) || response.Response.StatusCode != 200 {
                return false
        }
        lowerBody := strings.ToLower(string(response.Body))
        if len(lowerBody) < 10 {
                return false
        }

        // Vérification de la présence de la chaîne "repositoryformatversion"
        if strings.Contains(lowerBody, "repositoryformatversion") {
                event.Service.Software.Name = "GitConfigFile"
                event.Leak.Type = "config_leak"
                event.Leak.Severity = "high"
                event.AddTag("potential-git-config-leak")
                event.Summary = "Found repositoryformatversion in /static../.git/config:\n" + string(response.Body)
                return true
        }

        return false
}
