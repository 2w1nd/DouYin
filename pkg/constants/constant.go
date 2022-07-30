package constants

const (
	TablePre            = "douyin_"
	APIServiceName      = "api_service"
	UserServiceName     = "user_service"
	VideoServiceName    = "video_service"
	CommentServiceName  = "comment_service"
	FavoriteServiceName = "favorite_service"
	RelationServiceName = "relation_service"

	UserTableName          = TablePre + "user"
	VideoTableName         = TablePre + "video"
	RelationCountTableName = TablePre + "relation_count"
	FollowTableName        = TablePre + "follow"
	FollowerTableName      = TablePre + "follower"
	CommentTableName       = TablePre + "comment"
	CommentCountTableName  = TablePre + "comment_count"

	RelationFollowPre = "follow:"
	RelationFansPre   = "fans:"
	RelationCountPre  = "count:"
	FavoriteLikePre   = "like:"
	FavoriteVideoPre  = "video:"

	UserSalt     = "ByteDanceCamp"
	JWTSecretKey = "ByteDanceCamp3"
	IdentityKey  = "uid"

	QiNiuAccessKey = "keR1VefVxLVXyfcdg0E0KF4n8k72Ulcwc33fePrf"
	QiNiuSecretKey = "aKJAUNhALfgj1RqcIwBHd-513_o2yUV-wsh-qQdu"
	QiNiuBucket    = "micro-tiktok"
	QiNiuServer    = "http://data.mtt.dtpark.top/"

	EtcdAddress     = "47.96.157.235:2379"
	RedisAddress    = "127.0.0.1:16379"
	MySQLDefaultDSN = "root:5842020@tcp(1.117.93.124:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
)
