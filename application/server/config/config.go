//定义了一个配置模块，用于从 YAML 文件中读取配置信息，并将其存储到全局变量 GlobalConfig 中。
//配置信息包括服务器配置、Fabric 配置以及组织的详细配置。
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	Fabric FabricConfig `yaml:"fabric"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `yaml:"port"`
}

// FabricConfig Fabric配置
type FabricConfig struct {
	ChannelName   string                        `yaml:"channelName"`
	ChaincodeName string                        `yaml:"chaincodeName"`
	Organizations map[string]OrganizationConfig `yaml:"organizations"`
}

// OrganizationConfig 组织配置
type OrganizationConfig struct {
	MSPID        string `yaml:"mspID"`
	CertPath     string `yaml:"certPath"`
	KeyPath      string `yaml:"keyPath"`
	TLSCertPath  string `yaml:"tlsCertPath"`
	PeerEndpoint string `yaml:"peerEndpoint"`
	GatewayPeer  string `yaml:"gatewayPeer"`
}

var GlobalConfig Config

// InitConfig 初始化配置
func InitConfig() error {
	// 读取配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败：%v", err)
	}

	// 解析配置文件
	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("解析配置文件失败：%v", err)
	}

	return nil
}
