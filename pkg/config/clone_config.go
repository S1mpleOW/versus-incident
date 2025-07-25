package config

func cloneConfig(src *Config) *Config {
	if src == nil {
		return nil
	}

	// Create a new Config and copy all fields
	cloned := &Config{
		Name:       src.Name,
		Host:       src.Host,
		Port:       src.Port,
		PublicHost: src.PublicHost,
		Alert:      cloneAlertConfig(src.Alert),
		Queue:      cloneQueueConfig(src.Queue),
		OnCall:     cloneOnCallConfig(src.OnCall),
		Proxy:      cloneProxyConfig(src.Proxy),
		Redis:      cloneRedisConfig(src.Redis),
	}

	return cloned
}

// Helper function to deep clone the AlertConfig struct
func cloneAlertConfig(src AlertConfig) AlertConfig {
	return AlertConfig{
		DebugBody:  src.DebugBody,
		Slack:      cloneSlackConfig(src.Slack),
		Telegram:   cloneTelegramConfig(src.Telegram),
		Viber:      cloneViberConfig(src.Viber),
		Email:      cloneEmailConfig(src.Email),
		MSTeams:    cloneMSTeamsConfig(src.MSTeams),
		Lark:       cloneLarkConfig(src.Lark),
		GoogleChat: cloneGGChatConfig(src.GoogleChat),
	}
}

// Helper function to deep clone the SlackConfig struct
func cloneSlackConfig(src SlackConfig) SlackConfig {
	return SlackConfig{
		Enable:       src.Enable,
		Token:        src.Token,
		ChannelID:    src.ChannelID,
		TemplatePath: src.TemplatePath,
		MessageProperties: SlackMessageProperties{
			DisableButton: src.MessageProperties.DisableButton,
			ButtonText:    src.MessageProperties.ButtonText,
			ButtonStyle:   src.MessageProperties.ButtonStyle,
		},
	}
}

// Helper function to deep clone the TelegramConfig struct
func cloneTelegramConfig(src TelegramConfig) TelegramConfig {
	return TelegramConfig{
		Enable:       src.Enable,
		BotToken:     src.BotToken,
		ChatID:       src.ChatID,
		TemplatePath: src.TemplatePath,
		UseProxy:     src.UseProxy,
	}
}

// Helper function to deep clone the ViberConfig struct
func cloneViberConfig(src ViberConfig) ViberConfig {
	return ViberConfig{
		Enable:       src.Enable,
		APIType:      src.APIType,
		BotToken:     src.BotToken,
		UserID:       src.UserID,
		TemplatePath: src.TemplatePath,
		ChannelID:    src.ChannelID,
		UseProxy:     src.UseProxy,
	}
}

// Helper function to deep clone the EmailConfig struct
func cloneEmailConfig(src EmailConfig) EmailConfig {
	return EmailConfig{
		Enable:       src.Enable,
		SMTPHost:     src.SMTPHost,
		SMTPPort:     src.SMTPPort,
		Username:     src.Username,
		Password:     src.Password,
		To:           src.To,
		Subject:      src.Subject,
		TemplatePath: src.TemplatePath,
	}
}

// Helper function to deep clone the MSTeamsConfig struct
func cloneMSTeamsConfig(src MSTeamsConfig) MSTeamsConfig {
	// Create a copy of OtherPowerURLs map if it exists
	var otherPowerURLsCopy map[string]string
	if src.OtherPowerURLs != nil {
		otherPowerURLsCopy = make(map[string]string)
		for k, v := range src.OtherPowerURLs {
			otherPowerURLsCopy[k] = v
		}
	}

	return MSTeamsConfig{
		Enable:           src.Enable,
		TemplatePath:     src.TemplatePath,
		PowerAutomateURL: src.PowerAutomateURL,
		OtherPowerURLs:   otherPowerURLsCopy,
	}
}

// Helper function to deep clone the LarkConfig struct
func cloneLarkConfig(src LarkConfig) LarkConfig {
	// Create a copy of OtherWebhookURLs map if it exists
	var otherWebhookURLsCopy map[string]string
	if src.OtherWebhookURLs != nil {
		otherWebhookURLsCopy = make(map[string]string)
		for k, v := range src.OtherWebhookURLs {
			otherWebhookURLsCopy[k] = v
		}
	}

	return LarkConfig{
		Enable:           src.Enable,
		WebhookURL:       src.WebhookURL,
		TemplatePath:     src.TemplatePath,
		OtherWebhookURLs: otherWebhookURLsCopy,
		UseProxy:         src.UseProxy,
	}
}

// Helper function to deep clone the QueueConfig struct
func cloneQueueConfig(src QueueConfig) QueueConfig {
	return QueueConfig{
		Enable: src.Enable,
		SNS:    cloneSNSConfig(src.SNS),
		SQS:    cloneSQSConfig(src.SQS),
		PubSub: clonePubSubConfig(src.PubSub),
		AzBus:  cloneAzBusConfig(src.AzBus),
	}
}

// Helper function to deep clone the SNSConfig struct
func cloneSNSConfig(src SNSConfig) SNSConfig {
	return SNSConfig{
		Enable: src.Enable,
	}
}

// Helper function to deep clone the SQSConfig struct
func cloneSQSConfig(src SQSConfig) SQSConfig {
	return SQSConfig{
		Enable:   src.Enable,
		QueueURL: src.QueueURL,
	}
}

// Helper function to deep clone the PubSubConfig struct
func clonePubSubConfig(src PubSubConfig) PubSubConfig {
	return PubSubConfig{
		Enable: src.Enable,
	}
}

// Helper function to deep clone the AzBusConfig struct
func cloneAzBusConfig(src AzBusConfig) AzBusConfig {
	return AzBusConfig{
		Enable: src.Enable,
	}
}

// Helper function to deep clone the OnCallConfig struct
func cloneOnCallConfig(src OnCallConfig) OnCallConfig {
	return OnCallConfig{
		Enable:             src.Enable,
		InitializedOnly:    src.InitializedOnly,
		WaitMinutes:        src.WaitMinutes,
		Provider:           src.Provider,
		AwsIncidentManager: cloneAwsIncidentManagerConfig(src.AwsIncidentManager),
		PagerDuty:          clonePagerDutyConfig(src.PagerDuty),
	}
}

// Helper function to deep clone the AwsIncidentManagerConfig struct
func cloneAwsIncidentManagerConfig(src AwsIncidentManagerConfig) AwsIncidentManagerConfig {
	// Create a copy of OtherResponsePlanArns map if it exists
	var otherResponsePlanArnsCopy map[string]string
	if src.OtherResponsePlanArns != nil {
		otherResponsePlanArnsCopy = make(map[string]string)
		for k, v := range src.OtherResponsePlanArns {
			otherResponsePlanArnsCopy[k] = v
		}
	}

	return AwsIncidentManagerConfig{
		ResponsePlanArn:       src.ResponsePlanArn,
		OtherResponsePlanArns: otherResponsePlanArnsCopy,
	}
}

// Helper function to deep clone the PagerDutyConfig struct
func clonePagerDutyConfig(src PagerDutyConfig) PagerDutyConfig {
	// Create a copy of OtherRoutingKeys map if it exists
	var otherRoutingKeysCopy map[string]string
	if src.OtherRoutingKeys != nil {
		otherRoutingKeysCopy = make(map[string]string)
		for k, v := range src.OtherRoutingKeys {
			otherRoutingKeysCopy[k] = v
		}
	}

	return PagerDutyConfig{
		RoutingKey:       src.RoutingKey,
		OtherRoutingKeys: otherRoutingKeysCopy,
	}
}

// Helper function to deep clone the ProxyConfig struct
func cloneProxyConfig(src ProxyConfig) ProxyConfig {
	return ProxyConfig{
		URL:      src.URL,
		Username: src.Username,
		Password: src.Password,
	}
}

// Helper function to deep clone the RedisConfig struct
func cloneRedisConfig(src RedisConfig) RedisConfig {
	return RedisConfig{
		Host:               src.Host,
		Port:               src.Port,
		Password:           src.Password,
		DB:                 src.DB,
		InsecureSkipVerify: src.InsecureSkipVerify,
	}
}

func cloneOtherWebhookURLs(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}

	// Create a new map and copy the contents
	cloned := make(map[string]string, len(src))
	for k, v := range src {
		cloned[k] = v
	}
	return cloned
}

func cloneGGChatConfig(src GoogleChatConfig) GoogleChatConfig {
	return GoogleChatConfig{
		WebhookURL:       src.WebhookURL,
		TemplatePath:     src.TemplatePath,
		UseProxy:         src.UseProxy,
		OtherWebhookURLs: cloneOtherWebhookURLs(src.OtherWebhookURLs),
		Enable:           src.Enable,
	}
}
