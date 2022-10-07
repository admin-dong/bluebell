package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamVoteData struct {
	PostID    string `json:"post_id," binding:"required"`              //帖子id  //加string的目的是为了防止数字丢失帧
	Direction int8   `json:"direction,string" binding:"one of=1 0 -1"` //赞成票（1）还是反对票 （——1） 取消投票0
}

//ParamPostList 获取帖子列表query string 参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

type ParamCommunityPostList struct {
	ParamPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}

//投票数据
//
//type ParamVoteData struct {
//	//userID
//	PostID    int64 `json:"post_id,string"`   //加string的目的是为了防止数字丢失帧
//	Direction int   `json:"direction,string"` //赞成票（1）还是反对票 （——1）
//}
