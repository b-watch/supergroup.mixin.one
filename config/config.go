package config

import (
	"io/ioutil"
	"log"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const ConfigFile = "config.yaml"
const BuildVersion = "BUILD_VERSION"
const RouterModeHash = "hash"
const RouterModeHistory = "history"
const RouterModeHashSymbol = "/#"

type PaymentAsset struct {
	Symbol  string `yaml:"symbol" json:"symbol"`
	AssetId string `yaml:"asset_id" json:"asset_id"`
	Amount  string `yaml:"amount" json:"amount"`
}

type Shortcut struct {
	Icon    string `yaml:"icon" json:"icon"`
	LabelEn string `yaml:"label_en" json:"label_en"`
	LabelZh string `yaml:"label_zh" json:"label_zh"`
	Url     string `yaml:"url" json:"url"`
}

type ShortcutGroup struct {
	LabelEn string     `yaml:"label_en" json:"label_en"`
	LabelZh string     `yaml:"label_zh" json:"label_zh"`
	Items   []Shortcut `yaml:"shortcuts" json:"shortcuts"`
}

type Config struct {
	Service struct {
		Name             string `yaml:"name"`
		Environment      string `yaml:"enviroment"`
		HTTPListenPort   int    `yaml:"port"`
		HTTPResourceHost string `yaml:"host"`
		HTTPApiHost      string `yaml:"api_host"`
	} `yaml:"service"`
	Database struct {
		DatebaseUser     string `yaml:"username"`
		DatabasePassword string `yaml:"password"`
		DatabaseHost     string `yaml:"host"`
		DatabasePort     string `yaml:"port"`
		DatabaseName     string `yaml:"database_name"`
	} `yaml:"database"`
	System struct {
		RouterMode           string `yaml:"router_mode"`
		MessageShardModifier string `yaml:"message_shard_modifier"`
		MessageShardSize     int64  `yaml:"message_shard_size"`
		PriceAssetsEnable    bool   `yaml:"price_asset_enable"`
		ReadAssetsEnable     bool   `yaml:"read_assets_enable"`
		AudioMessageEnable   bool   `yaml:"audio_message_enable"`
		ImageMessageEnable   bool   `yaml:"image_message_enable"`
		VideoMessageEnable   bool   `yaml:"video_message_enable"`
		ContactMessageEnable bool   `yaml:"contact_message_enable"`

		RewardsEnable                   bool     `yaml:"rewards_enable"`
		RewardsMinAmountBase            string   `yaml:"rewards_min_amount_base"`
		RewardsAssetList                []string `yaml:"rewards_asset_list"`
		RedPacketMinAmountBase          string   `yaml:"redpacket_min_amount_base"`
		RedPacketAssetList              []string `yaml:"redpacket_asset_list"`
		RedPacketNormDistSigmaMeanRatio string   `yaml:"redpacket_normal_distribution_sigma_mean_ratio"`

		LimitMessageFrequency    bool     `yaml:"limit_message_frequency"`
		OperatorList             []string `yaml:"operator_list"`
		Operators                map[string]bool
		DetectQRCodeEnabled      bool           `yaml:"detect_image"`
		DetectLinkEnabled        bool           `yaml:"detect_link"`
		SensitiveWords           string         `yaml:"sensitive_words"`
		ProhibitedMessageEnabled bool           `yaml:"prohibited_message"`
		PaymentAssetId           string         `yaml:"payment_asset_id"`
		PaymentAmount            string         `yaml:"payment_amount"`
		PayToJoin                bool           `yaml:"pay_to_join"`
		InviteToJoin             bool           `yaml:"invite_to_join"`
		AutoEstimate             bool           `yaml:"auto_estimate"`
		AutoEstimateCurrency     string         `yaml:"auto_estimate_currency"`
		AutoEstimateBase         string         `yaml:"auto_estimate_base"`
		AccpetPaymentAssetList   []PaymentAsset `yaml:"accept_asset_list"`
		AccpetWeChatPayment      bool           `yaml:"accept_wechat_payment"`
		WeChatPaymentAmount      string         `yaml:"wechat_payment_amount"`
		AccpetCouponPayment      bool           `yaml:"accept_coupon_payment"`
	} `yaml:"system"`
	Appearance struct {
		HomeWelcomeMessage string          `yaml:"home_welcome_message"`
		HomeShortcutGroups []ShortcutGroup `yaml:"home_shortcut_groups"`
	} `yaml:"appearance"`
	MessageTemplate struct {
		WelcomeMessage                string `yaml:"welcome_message"`
		GroupRedPacket                string `yaml:"group_redpacket"`
		GroupRedPacketShortDesc       string `yaml:"group_redpacket_short_desc"`
		GroupRedPacketDesc            string `yaml:"group_redpacket_desc"`
		GroupOpenedRedPacket          string `yaml:"group_opened_redpacket"`
		MessageTipsGuest              string `yaml:"message_tips_guest"`
		MessageProhibit               string `yaml:"message_prohibit"`
		MessageAllow                  string `yaml:"message_allow"`
		MessageTipsJoin               string `yaml:"message_tips_join"`
		MessageTipsJoinUser           string `yaml:"message_tips_join_user"`
		MessageTipsJoinUserProhibited string `yaml:"message_tips_join_user_prohibited"`
		MessageTipsHelp               string `yaml:"message_tips_help"`
		MessageTipsHelpBtn            string `yaml:"message_tips_help_btn"`
		MessageTipsUnsubscribe        string `yaml:"message_tips_unsubscribe"`
		MessageTipsTooMany            string `yaml:"message_tips_too_many"`
		MessageTipsRewards            string `yaml:"message_tips_rewards"`
		MessageCommandsInfo           string `yaml:"message_commands_info"`
		MessageCommandsInfoResp       string `yaml:"message_commands_info_resp"`
	} `yaml:"message_template"`
	Wechat struct {
		AppId          string `yaml:"app_id"`
		AppSecret      string `yaml:"app_secret"`
		Token          string `yaml:"token"`
		EncodingAESKey string `yaml:"encodine_aes_key"`
		MchId          string `yaml:"mch_id"`
		MchKey         string `yaml:"mch_key"`
		NotifyUrl      string `yaml:"notify_url"`
	} `yaml:"wechat"`
	Mixin struct {
		ClientId        string `yaml:"client_id"`
		ClientSecret    string `yaml:"client_secret"`
		SessionAssetPIN string `yaml:"session_asset_pin"`
		PinToken        string `yaml:"pin_token"`
		SessionId       string `yaml:"session_id"`
		SessionKey      string `yaml:"session_key"`
	} `yaml:"mixin"`
	Plugins []struct {
		SharedLibrary string                 `yaml:"shared_library"`
		Config        map[string]interface{} `yaml:"config"`
	}
}

type ExportedConfig struct {
	MixinClientId          string          `json:"mixin_client_id"`
	HTTPResourceHost       string          `json:"host"`
	HTTPApiHost            string          `json:"api_host"`
	AutoEstimate           bool            `json:"auto_estimate"`
	AutoEstimateCurrency   string          `json:"auto_estimate_currency"`
	AutoEstimateBase       string          `json:"auto_estimate_base"`
	AccpetPaymentAssetList []PaymentAsset  `json:"accept_asset_list"`
	AccpetWeChatPayment    bool            `json:"accept_wechat_payment"`
	WeChatPaymentAmount    string          `json:"wechat_payment_amount"`
	AccpetCouponPayment    bool            `json:"accept_coupon_payment"`
	PayToJoin              bool            `json:"pay_to_join"`
	InviteToJoin           bool            `json:"invite_to_join"`
	RewardsEnable          bool            `json:"rewards_enable"`
	RewardsMinAmountBase   string          `json:"rewards_min_amount_base"`
	RedPacketMinAmountBase string          `json:"redpacket_min_amount_base"`
	HomeWelcomeMessage     string          `json:"home_welcome_message"`
	HomeShortcutGroups     []ShortcutGroup `json:"home_shortcut_groups"`
	ServiceName            string          `json:"service_name"`
}

var AppConfig *Config

func LoadConfig(dir string) {
	data, err := ioutil.ReadFile(path.Join(dir, ConfigFile))
	if err != nil {
		log.Panicln(err)
	}
	AppConfig = &Config{}
	err = yaml.Unmarshal(data, AppConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	AppConfig.System.Operators = make(map[string]bool)
	for _, op := range AppConfig.System.OperatorList {
		AppConfig.System.Operators[op] = true
	}
}

func GetExported() ExportedConfig {
	var exc ExportedConfig
	exc.MixinClientId = AppConfig.Mixin.ClientId
	exc.HTTPResourceHost = AppConfig.Service.HTTPResourceHost
	exc.HTTPApiHost = AppConfig.Service.HTTPApiHost
	exc.AutoEstimate = AppConfig.System.AutoEstimate
	exc.AutoEstimateCurrency = AppConfig.System.AutoEstimateCurrency
	exc.AutoEstimateBase = AppConfig.System.AutoEstimateBase
	exc.AccpetPaymentAssetList = AppConfig.System.AccpetPaymentAssetList
	exc.AccpetWeChatPayment = AppConfig.System.AccpetWeChatPayment
	exc.WeChatPaymentAmount = AppConfig.System.WeChatPaymentAmount
	exc.AccpetCouponPayment = AppConfig.System.AccpetCouponPayment
	exc.PayToJoin = AppConfig.System.PayToJoin
	exc.InviteToJoin = AppConfig.System.InviteToJoin
	exc.RewardsEnable = AppConfig.System.RewardsEnable
	exc.RewardsMinAmountBase = AppConfig.System.RewardsMinAmountBase
	exc.RedPacketMinAmountBase = AppConfig.System.RedPacketMinAmountBase
	exc.HomeWelcomeMessage = AppConfig.Appearance.HomeWelcomeMessage
	exc.HomeShortcutGroups = AppConfig.Appearance.HomeShortcutGroups
	exc.ServiceName = AppConfig.Service.Name
	return exc
}
