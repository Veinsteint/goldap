package response

import "goldap-server/model"

type PendingUserListRsp struct {
	Total        int                 `json:"total"`
	PendingUsers []*model.PendingUser `json:"pendingUsers"`
}

