package constant

import "github.com/openimsdk/protocol/constant"

const (
	MountConfigFilePath = "CONFIG_PATH"
	KUBERNETES          = "kubernetes"
	ETCD                = "etcd"
)

const (
	// verificationCode used for.
	VerificationCodeForRegister      = 1 // Register
	VerificationCodeForResetPassword = 2 // Reset password
	VerificationCodeForLogin         = 3 // Login
)

const LogFileName = "chat.log"

// block unblock.
const (
	BlockUser   = 1
	UnblockUser = 2
)

// AccountType.
const (
	Email   = "email"
	Phone   = "phone"
	Account = "account"
)

// Mode.
const (
	UserMode  = "user"
	AdminMode = "admin"
)

const DefaultAdminLevel = 100

// user level.
const (
	NormalAdmin       = 80
	AdvancedUserLevel = 100
)

// AddFriendCtrl.
const (
	OrdinaryUserAddFriendEnable  = 1  // Allow ordinary users to add friends
	OrdinaryUserAddFriendDisable = -1 // Do not allow ordinary users to add friends
)

const (
	NormalUser = 1
	AdminUser  = 2
)

// mini-app
const (
	StatusOnShelf = 1 // OnShelf
	StatusUnShelf = 2 // UnShelf
)

const (
	LimitNil             = 0 // None
	LimitEmpty           = 1 // Neither are restricted
	LimitOnlyLoginIP     = 2 // Only login is restricted
	LimitOnlyRegisterIP  = 3 // Only registration is restricted
	LimitLoginIP         = 4 // Restrict login
	LimitRegisterIP      = 5 // Restrict registration
	LimitLoginRegisterIP = 6 // Restrict both login and registration
)

const (
	InvitationCodeAll    = 0 // All
	InvitationCodeUsed   = 1 // Used
	InvitationCodeUnused = 2 // Unused
)

const (
	RpcOpUserID   = constant.OpUserID
	RpcOpUserType = "opUserType"
)

const RpcCustomHeader = constant.RpcCustomHeader

const NeedInvitationCodeRegisterConfigKey = "needInvitationCodeRegister"

const (
	DefaultAllowVibration = 1
	DefaultAllowBeep      = 1
	DefaultAllowAddFriend = 1
)

const (
	FinDAllUser    = 0
	FindNormalUser = 1
)

const CtxApiToken = "api-token"

const (
	AccountRegister = iota
	EmailRegister
	PhoneRegister
)

const (
	GenderFemale  = 0 // female
	GenderMale    = 1 // male
	GenderUnknown = 2 // unknown
)

// Credential Type
const (
	CredentialAccount = iota
	CredentialPhone
	CredentialEmail
)

const (
	Follow    = 0
	Subscribe = 1
	Reply     = 2
	Like      = 3
	Collect   = 4
)

const (
	NotAllow = 0
	Allow    = 1
)

const (
	NotLiked = 0
	Liked    = 1
)

const (
	NotCollected = 0
	Collected    = 1
)

const (
	NotForwarded = 0
	Forwarded    = 1
)

const (
	NotCommented = 0
	Commented    = 1
)

const (
	NotFollowed = 0
	Followed    = 1
)

const (
	NotSubscribed = 0
	Subscribed    = 1
)

const (
	PostMediaTypePicture = 0
	PostMediaTypeVideo   = 1
)

const (
	Pinned   = 1
	UnPinned = 0
)

const (
	SendPrivateRedPacket   = 1000 //聊发送红包
	SendLuckRedPacket      = 1001 //群聊拼手气红包
	SendExclusiveRedPacket = 1002 //群聊用户专属红包

	ReceivePrivateRedPacket   = 1003 //领取聊发送红包
	ReceiveLuckRedPacket      = 1004 //领取群聊拼手气红包
	ReceiveExclusiveRedPacket = 1005 //领取群聊用户专属红包

	RefundRedPacket = 1005 //退款
)

const (
	StealthUser    = 1
	NotStealthUser = 0
)
