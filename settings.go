package gbcobservehttp

import coreenv "goark.dev/goark/core/env"

type settings struct{ enabled, serverEnabled, clientEnabled *bool }

func newSettings(environment coreenv.Environment, options []Option) (settings, error) {
	resolved := settings{}
	if environment != nil {
		for key, target := range map[string]**bool{PropertyEnabled: &resolved.enabled, PropertyServerEnabled: &resolved.serverEnabled, PropertyClientEnabled: &resolved.clientEnabled} {
			value, err := coreenv.ResolveValueAs[bool](environment, "${"+key+":true}")
			if err != nil {
				return settings{}, err
			}
			copied := value
			*target = &copied
		}
	}
	for _, option := range options {
		if option != nil {
			option(&resolved)
		}
	}
	if resolved.enabled == nil {
		value := DefaultEnabled
		resolved.enabled = &value
	}
	if resolved.serverEnabled == nil {
		value := true
		resolved.serverEnabled = &value
	}
	if resolved.clientEnabled == nil {
		value := true
		resolved.clientEnabled = &value
	}
	if !*resolved.enabled {
		disabled := false
		resolved.serverEnabled = &disabled
		resolved.clientEnabled = &disabled
	}
	return resolved, nil
}
