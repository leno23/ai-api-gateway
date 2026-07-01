package model

import (
	"time"

	"github.com/lib/pq"
)

const (
	RoleUser       int16 = 1
	RoleAdmin      int16 = 10
	RoleSuperAdmin int16 = 100

	UserStatusActive   int16 = 1
	UserStatusDisabled int16 = 0

	APIKeyStatusActive   int16 = 1
	APIKeyStatusDisabled int16 = 0

	ChannelStatusActive int16 = 1

	// Redeem: migration default status = 1 for available rows
	RedeemStatusAvailable int16 = 1
	RedeemStatusUsed      int16 = 2

	QuotaLogTypeConsume       int16 = 1
	QuotaLogTypeRecharge      int16 = 2
	QuotaLogTypeRefund        int16 = 3
	QuotaLogTypeAdjust        int16 = 4
	QuotaLogTypeInvite        int16 = 5
	QuotaLogTypeConsumeRebate int16 = 6

	RebateTaskPending int16 = 0
	RebateTaskDone    int16 = 1
	RebateTaskFailed  int16 = 2
)

type User struct {
	ID            int64     `gorm:"primaryKey"`
	Username      string    `gorm:"uniqueIndex;size:64;not null"`
	Email         string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash  string    `gorm:"column:password_hash;size:255;not null"`
	Role          int16     `gorm:"default:1;not null"`
	Status        int16     `gorm:"default:1;not null"`
	Balance       float64   `gorm:"type:decimal(16,6);default:0"`
	Quota         int64     `gorm:"default:0;not null"`
	UsedQuota     int64     `gorm:"column:used_quota;default:0;not null"`
	InviteCode    string    `gorm:"column:invite_code;uniqueIndex;size:32"`
	InvitedBy     *int64    `gorm:"column:invited_by"`
	TokenGroupID       *int64    `gorm:"column:token_group_id"`
	AffiliatePending   int64     `gorm:"column:affiliate_pending;default:0;not null"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

type TokenGroup struct {
	ID         int64     `gorm:"primaryKey"`
	Slug       string    `gorm:"uniqueIndex;size:64;not null"`
	Name       string    `gorm:"size:128;not null"`
	Multiplier float64   `gorm:"type:decimal(6,2);default:1;not null"`
	Status     int16     `gorm:"default:1;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (TokenGroup) TableName() string { return "token_groups" }

func (User) TableName() string { return "users" }

type APIKey struct {
	ID            int64          `gorm:"primaryKey"`
	UserID        int64          `gorm:"index:idx_api_keys_user_id;not null"`
	Name          string         `gorm:"size:128;not null"`
	KeyHash       string         `gorm:"column:key_hash;uniqueIndex:idx_api_keys_key_hash;size:255;not null"`
	KeyPrefix     string         `gorm:"column:key_prefix;size:16;not null"`
	Status        int16          `gorm:"default:1;not null"`
	TokenGroupID  *int64         `gorm:"column:token_group_id"`
	QuotaLimit    *int64         `gorm:"column:quota_limit"`
	UsedQuota     int64          `gorm:"column:used_quota;default:0;not null"`
	Models        pq.StringArray `gorm:"type:text[]"`
	IPWhitelist   pq.StringArray `gorm:"column:ip_whitelist;type:text[]"`
	RateLimit     int            `gorm:"column:rate_limit;default:60"`
	ExpiresAt     *time.Time     `gorm:"column:expires_at"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
}

func (APIKey) TableName() string { return "api_keys" }

type Channel struct {
	ID            int64          `gorm:"primaryKey"`
	Name          string         `gorm:"size:128;not null"`
	Provider      string         `gorm:"size:32;not null"`
	BaseURL       string         `gorm:"column:base_url;size:512;not null"`
	APIKey        string         `gorm:"column:api_key;size:512;not null"`
	Models        pq.StringArray `gorm:"type:text[];default:'{}'"`
	ModelMapping  []byte         `gorm:"column:model_mapping;type:jsonb"`
	Priority      int            `gorm:"default:0;not null"`
	Weight        int            `gorm:"default:1;not null"`
	Status        int16          `gorm:"default:1;not null"`
	MaxConcurrent int            `gorm:"column:max_concurrent;default:100"`
	RateLimit     int            `gorm:"column:rate_limit;default:1000"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
}

func (Channel) TableName() string { return "channels" }

type ModelPrice struct {
	ID              int64          `gorm:"primaryKey"`
	Model           string         `gorm:"column:model;uniqueIndex;size:128;not null"`
	DisplayName     string         `gorm:"column:display_name;size:128"`
	Provider        string         `gorm:"size:32;default:openai;not null"`
	EndpointType    string         `gorm:"column:endpoint_type;size:32;default:openai;not null"`
	PromptPrice     int64          `gorm:"column:prompt_price;not null"`
	CompletionPrice int64          `gorm:"column:completion_price;not null"`
	CacheReadPrice  int64          `gorm:"column:cache_read_price;not null"`
	CacheWritePrice int64          `gorm:"column:cache_write_price;not null"`
	UnitPrice       int64          `gorm:"column:unit_price;not null"`
	BillingType     int16          `gorm:"column:billing_type;default:1;not null"`
	Tags            pq.StringArray `gorm:"type:text[]"`
	Currency        string         `gorm:"size:8;default:CNY;not null"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
}

func (ModelPrice) TableName() string { return "model_prices" }

type RequestLog struct {
	ID               int64     `gorm:"primaryKey"`
	UserID           int64     `gorm:"index:idx_request_logs_user_id;not null"`
	APIKeyID         *int64    `gorm:"column:api_key_id"`
	ChannelID        *int64    `gorm:"column:channel_id"`
	RequestID        string    `gorm:"column:request_id;size:64"`
	TokenName        string    `gorm:"column:token_name;size:128"`
	TokenGroup       string    `gorm:"column:token_group;size:64"`
	Model            string    `gorm:"size:128"`
	RequestMethod    string    `gorm:"column:request_method;size:16"`
	RequestPath      string    `gorm:"column:request_path;size:256"`
	PromptTokens     int       `gorm:"column:prompt_tokens"`
	CompletionTokens int       `gorm:"column:completion_tokens"`
	TotalTokens      int       `gorm:"column:total_tokens"`
	CostQuota        int64     `gorm:"column:cost_quota"`
	LatencyMs        *int      `gorm:"column:latency_ms"`
	TimeToFirstMs    *int      `gorm:"column:time_to_first_ms"`
	BillingDetail    []byte    `gorm:"column:billing_detail;type:jsonb"`
	StatusCode       *int      `gorm:"column:status_code"`
	ErrorMessage     string    `gorm:"column:error_message;type:text"`
	CreatedAt        time.Time `gorm:"autoCreateTime;index:idx_request_logs_user_id,priority:2"`
}

func (RequestLog) TableName() string { return "request_logs" }

type TaskLog struct {
	ID          int64      `gorm:"primaryKey"`
	UserID      int64      `gorm:"index:idx_task_logs_user_submitted;not null"`
	TaskID      string     `gorm:"column:task_id;size:128;not null"`
	Platform    string     `gorm:"size:64"`
	TaskType    string     `gorm:"column:task_type;size:64"`
	Status      string     `gorm:"size:32;not null"`
	Progress    int        `gorm:"default:0;not null"`
	Detail      string     `gorm:"type:text"`
	SubmittedAt time.Time  `gorm:"column:submitted_at;not null"`
	FinishedAt  *time.Time `gorm:"column:finished_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
}

func (TaskLog) TableName() string { return "task_logs" }

type Announcement struct {
	ID        int64      `gorm:"primaryKey"`
	Title     string     `gorm:"size:256;not null"`
	Content   string     `gorm:"type:text;not null"`
	Level     string     `gorm:"size:32;default:info;not null"`
	Placement string     `gorm:"size:64;default:home;not null"`
	Status    int16      `gorm:"default:1;not null"`
	StartsAt  *time.Time `gorm:"column:starts_at"`
	EndsAt    *time.Time `gorm:"column:ends_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

func (Announcement) TableName() string { return "announcements" }

type RechargeRecord struct {
	ID            int64      `gorm:"primaryKey"`
	UserID        int64      `gorm:"index;not null"`
	Amount        float64    `gorm:"type:decimal(16,2);not null"`
	PaymentMethod string     `gorm:"column:payment_method;size:32"`
	TradeNo       string     `gorm:"column:trade_no;uniqueIndex;size:128"`
	Status        int16      `gorm:"default:0;not null"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	CompletedAt   *time.Time `gorm:"column:completed_at"`
}

func (RechargeRecord) TableName() string { return "recharge_records" }

type RedeemCode struct {
	ID        int64      `gorm:"primaryKey"`
	Code      string     `gorm:"uniqueIndex;size:64;not null"`
	Quota     int64      `gorm:"not null"`
	UsedBy    *int64     `gorm:"column:used_by"`
	Status    int16      `gorm:"default:1;not null"`
	ExpiresAt *time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UsedAt    *time.Time `gorm:"column:used_at"`
}

func (RedeemCode) TableName() string { return "redeem_codes" }

type QuotaLog struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"index:idx_quota_logs_user;priority:1;not null"`
	Delta     int64     `gorm:"not null"`
	Balance   int64     `gorm:"not null"`
	Type      int16     `gorm:"not null"`
	Reference string    `gorm:"size:128"`
	Remark    string    `gorm:"size:512"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_quota_logs_user,priority:2"`
}

func (QuotaLog) TableName() string { return "quota_logs" }

type InviteRecord struct {
	ID           int64     `gorm:"primaryKey"`
	InviterID    int64     `gorm:"column:inviter_id;index:idx_invite_records_inviter;not null"`
	InviteeID    int64     `gorm:"column:invitee_id;uniqueIndex;not null"`
	InviterBonus int64     `gorm:"column:inviter_bonus;not null"`
	InviteeBonus int64     `gorm:"column:invitee_bonus;not null"`
	Status       int16     `gorm:"default:1;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (InviteRecord) TableName() string { return "invite_records" }

// RebateTask is an optional post-consume rebate credit (see REBATE_* config).
type RebateTask struct {
	ID          int64      `gorm:"primaryKey"`
	UserID      int64      `gorm:"index;not null"`
	Amount      int64      `gorm:"not null"`
	Status      int16      `gorm:"default:0;not null"`
	Ref         string     `gorm:"size:128"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	ProcessedAt *time.Time `gorm:"column:processed_at"`
}

func (RebateTask) TableName() string { return "rebate_tasks" }
