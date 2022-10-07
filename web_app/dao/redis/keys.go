package redis

//redis key

//redis key 尽量使用命名空间的方式 区分不同的key  方便业务拆分
const (
	Prefix             = "bluebell:"
	KeyPostTimeZset    = "post:time"   //zset;帖子以发帖时间
	KeyPostScoreZset   = "post:score"  //zset; 帖子以及投票的分数
	KeyPostVotedZsetPf = "post:voted:" //zset 记录以及投票的类型  不完整的 参数是post——id
)

//给redis key增加前缀
func getReadisKey(key string) string {
	return Prefix + key
}
