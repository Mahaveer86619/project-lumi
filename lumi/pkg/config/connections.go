package config

type ConnectionsConfig struct {
	WhatsappConnectionStatus string
}

var GConnectionsConfig *ConnectionsConfig

func InitConnectionsConfig() {
	GConnectionsConfig = &ConnectionsConfig{
		WhatsappConnectionStatus: "unknown",
	}
}

func SetWhatsappConnectionStatus(status string) {
	GConnectionsConfig.WhatsappConnectionStatus = status
}

func GetWhatsappConnectionStatus() string {
	return GConnectionsConfig.WhatsappConnectionStatus
}
