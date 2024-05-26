package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/services"
	"github.com/gandarfh/httui/pkg/common"
)

func Connect(program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		organization, err := services.OrganizationShow()
		if err != nil {
			log.Println("failed to get organization value:", err.Error())
			return nil
		}

		// Load client cert
		cert, err := tls.X509KeyPair([]byte(organization.MqttCertPEM), []byte(organization.MqttCertKey))
		if err != nil {
			log.Fatalf("failed to parse client certificate: %v", err)
		}

		// Load CA cert
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM([]byte(organization.MqttCertCA)); !ok {
			log.Fatalf("failed to parse root certificate")
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}

		opts := MQTT.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("ssl://%s:8883", organization.MqttEndpoint))
		opts.SetClientID(config.Config.Settings.DeviceID)
		opts.SetTLSConfig(tlsConfig)

		messageHandler := func(client MQTT.Client, msg MQTT.Message) {
			sync := common.Sync{}
			json.Unmarshal(msg.Payload(), &sync)
			program.Send(sync)
		}

		opts.OnConnect = func(c MQTT.Client) {
			token := c.Subscribe(organization.MqttTopic, 0, messageHandler)

			token.Wait()
			if token.Error() != nil {
				log.Fatalf("subscribe error: %v", token.Error())
			}
		}

		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalf("failed to connect to MQTT broker: %v", token.Error())
		}

		log.Println("MQTT Connected!", organization.MqttTopic)

		return nil
	}
}
