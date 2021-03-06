package splunk

import (
	"terraform-provider-splunk/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type SplunkProvider struct {
	Client *client.Client
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: providerDataSources(),
		ResourcesMap:   providerResources(),
		ConfigureFunc:  providerConfigure,
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_URL", "localhost:8089"),
			Description: "Splunk instance URL",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_USERNAME", "admin"),
			Description: "Splunk instance admin username",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_PASSWORD", "changeme"),
			Description: "Splunk instance password",
		},
		"insecure_skip_verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_INSECURE_SKIP_VERIFY", true),
			Description: "insecure skip verification flag",
		},
	}
}

// Returns a map of splunk resources for configuration
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"splunk_apps_local":                  appsLocal(),
		"splunk_authentication_users":        authenticationUsers(),
		"splunk_authorization_roles":         authorizationRoles(),
		"splunk_global_http_event_collector": globalHttpEventCollector(),
		"splunk_inputs_http_event_collector": inputsHttpEventCollector(),
		"splunk_inputs_script":               inputsScript(),
		"splunk_inputs_monitor":              inputsMonitor(),
		"splunk_inputs_udp":                  inputsUDP(),
		"splunk_inputs_tcp_raw":              inputsTCPRaw(),
		"splunk_inputs_tcp_cooked":           inputsTCPCooked(),
		"splunk_inputs_tcp_splunk_tcp_token": inputsTCPSplunkTCPToken(),
		"splunk_inputs_tcp_ssl":              inputsTCPSSL(),
		"splunk_outputs_tcp_default":         outputsTCPDefault(),
		"splunk_outputs_tcp_server":          outputsTCPServer(),
		"splunk_outputs_tcp_group":           outputsTCPGroup(),
		"splunk_outputs_tcp_syslog":          outputsTCPSyslog(),
		"splunk_saved_searches":              savedSearches(),
		"splunk_indexes":                     index(),
		"splunk_configs_conf":                configsConf(),
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise splunk client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := client.NewSplunkdClient("", [2]string{d.Get("username").(string), d.Get("password").(string)},
		d.Get("url").(string), client.NewSplunkdHTTPClient(time.Second*30, d.Get("insecure_skip_verify").(bool)))
	err := client.Login()

	if err != nil {
		return client, err
	}

	provider := &SplunkProvider{
		Client: client,
	}

	return provider, nil
}
