package notifications

import (
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/eventsystem"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/configstore"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"strconv"
)

type Plugin struct{}

func RegisterPlugin() {
	plugin := &Plugin{}
	common.RegisterPlugin(plugin)

	common.GORM.AutoMigrate(&Config{})
	configstore.RegisterConfig(configstore.SQL, &Config{})

}

func (p *Plugin) Name() string {
	return "Notifications"
}

func (p *Plugin) InitBot() {
	eventsystem.AddHandler(bot.RedisWrapper(HandleGuildMemberAdd), eventsystem.EventGuildMemberAdd)
	eventsystem.AddHandler(bot.RedisWrapper(HandleGuildMemberRemove), eventsystem.EventGuildMemberRemove)
	eventsystem.AddHandlerBefore(HandleChannelUpdate, eventsystem.EventChannelUpdate, bot.StateHandlerPtr)
}

type Config struct {
	configstore.GuildConfigModel
	JoinServerEnabled bool   `json:"join_server_enabled" schema:"join_server_enabled"`
	JoinServerChannel string `json:"join_server_channel" schema:"join_server_channel" valid:"channel,true"`
	JoinServerMsg     string `json:"join_server_msg" schema:"join_server_msg" valid:"template,2000"`

	JoinDMEnabled bool   `json:"join_dm_enabled" schema:"join_dm_enabled"`
	JoinDMMsg     string `json:"join_dm_msg" schema:"join_dm_msg" valid:"template,2000"`

	LeaveEnabled bool   `json:"leave_enabled" schema:"leave_enabled"`
	LeaveChannel string `json:"leave_channel" schema:"leave_channel" valid:"channel,true"`
	LeaveMsg     string `json:"leave_msg" schema:"leave_msg" valid:"template,500"`

	TopicEnabled bool   `json:"topic_enabled" schema:"topic_enabled"`
	TopicChannel string `json:"topic_channel" schema:"topic_channel" valid:"channel,true"`
}

func (c *Config) JoinServerChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.JoinServerChannel, 10, 64)
	return
}

func (c *Config) LeaveChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.LeaveChannel, 10, 64)
	return
}

func (c *Config) TopicChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.TopicChannel, 10, 64)
	return
}

func (c *Config) GetName() string {
	return "general_notifications"
}

func (c *Config) TableName() string {
	return "general_notification_configs"
}

var DefaultConfig = &Config{}

func GetConfig(guildID int64) *Config {
	var conf Config
	err := configstore.Cached.GetGuildConfig(context.Background(), guildID, &conf)
	if err != nil {
		if err != configstore.ErrNotFound {
			log.WithError(err).Error("Failed retrieving config")
		}
		return &Config{
			JoinServerMsg: "<@{{.User.ID}}> Joined!",
			LeaveMsg:      "**{{.User.Username}}** Left... :'(",
		}
	}
	return &conf
}
