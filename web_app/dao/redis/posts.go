package redis

import (
	"web_app/models"

	"github.com/go-redis/redis"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带order参数确定要查询的order key
	key := getReadisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getReadisKey(KeyPostScoreZset)
	}
	//2.确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	//3.ZREVRANGE 按分数从大到小的顺序指定查询数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

// 根据ids 每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getReadisKey(KeyPostVotedZsetPf + id)
	//	//查找key中分数是1的元素的数量>统计每篇帖子的赞成票的数量
	//	v := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用pipline一次发送多条命令，较少rtt
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getReadisKey(KeyPostVotedZsetPf + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
