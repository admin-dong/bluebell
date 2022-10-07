package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

//1.用户投票的数据

//本项目是使用简化版的投票分数
//投一票就加432分  86400/200 ——> 需要200张赞成票可以给你的 帖子续一天

/*/投票的几种情况
//direction=1 时候,有两种情况  --> 更新分数和投票记录
//1.之前没有投过票，现在投赞成票  --> 更新分数和投票记录     差值的绝对值 1  +432
//2.之前投反对票,现在改投赞成票    --> 更新分数和投票记录    差值的绝对值2   +432*2
//direction=0 时候，有两种情况   --> 更新分数和投票记录
1.之前投过反对票。现在要取消投票  --> 更新分数和投票记录   差值的绝对值1   +432
2.之前投过赞成票。现在要取消投票  --> 更新分数和投票记录  差值的绝对值1   -432


direcion=-1时候，有两种情况：
 1.之前没投过票。现在要投反对票   --> 更新分数和投票记录  差值的绝对值  1  -432
2.之前投赞成票，现在改投反对票  --> 更新分数和投票记录   差值的绝对值 2   -432*2


投票的限制：
  每个帖子自发表之日起一个星期之允许用户投票，超过一个星期就不允许在投票了
1.到期之后 将redis中报错的赞成票数以及反对票数存储到mysql表中
2.到期之后删除那个 KeyPostVotedZsetPf（记录以及投票的类型  不完整的 参数是post——id）



*/
const (
	oneWeekSeconds = 7 * 24 * 3600
	scorePerVote   = 432 //每一票值多少
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("不允许重复投票")
)

func CreatePost(postID int64) error {
	//帖子时间
	pipelne := rdb.TxPipeline()
	pipelne.ZAdd(getReadisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数
	pipelne.ZAdd(getReadisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipelne.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票的限制
	//去redis取帖子发布时间
	postTime := rdb.ZScore(getReadisKey(KeyPostTimeZset), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekSeconds {
		return ErrVoteTimeExpire
	}
	//2.更新帖子分数

	//2和3 要放到一个pipline
	//先查当前用户给当前帖子之前的 投票记录
	ov := rdb.ZScore(getReadisKey(KeyPostVotedZsetPf+postID), userID).Val()
	//如果这一次投票的值和之前保存的值有一致的话，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepested
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算2次投票的差值
	pipelne := rdb.TxPipeline()
	pipelne.ZIncrBy(getReadisKey(KeyPostScoreZset), op*diff*scorePerVote, postID)
	//3.记录与用户为该帖子投过票
	if value == 0 {
		pipelne.ZRem(getReadisKey(KeyPostVotedZsetPf+postID), postID)
	} else {
		pipelne.ZAdd(getReadisKey(KeyPostVotedZsetPf+postID), redis.Z{
			Score:  value, //当前的用户投票赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipelne.Exec()
	return err
}
