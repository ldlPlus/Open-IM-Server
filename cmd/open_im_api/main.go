package main

import (
	apiAuth "Open_IM/internal/api/auth"
	apiChat "Open_IM/internal/api/chat"
	"Open_IM/internal/api/conversation"
	"Open_IM/internal/api/friend"
	"Open_IM/internal/api/group"
	"Open_IM/internal/api/manage"
	apiThird "Open_IM/internal/api/third"
	"Open_IM/internal/api/user"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
	//"syscall"
)

func main() {
	r := gin.Default()
	r.Use(utils.CorsHandler())
	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/update_user_info", user.UpdateUserInfo)
		userRouterGroup.POST("/get_user_info", user.GetUserInfo)
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		//	friendRouterGroup.POST("/get_friends_info", friend.GetFriendsInfo)
		friendRouterGroup.POST("/add_friend", friend.AddFriend)                     //1
		friendRouterGroup.POST("/get_friend_apply_list", friend.GetFriendApplyList) //1
		friendRouterGroup.POST("/get_self_apply_list", friend.GetSelfApplyList)     //1
		friendRouterGroup.POST("/get_friend_list", friend.GetFriendList)
		friendRouterGroup.POST("/add_blacklist", friend.AddBlacklist)
		friendRouterGroup.POST("/get_blacklist", friend.GetBlacklist)
		friendRouterGroup.POST("/remove_blacklist", friend.RemoveBlacklist)
		friendRouterGroup.POST("/delete_friend", friend.DeleteFriend)
		friendRouterGroup.POST("/add_friend_response", friend.AddFriendResponse) //1
		friendRouterGroup.POST("/set_friend_comment", friend.SetFriendComment)
		friendRouterGroup.POST("/is_friend", friend.IsFriend)
		friendRouterGroup.POST("/import_friend", friend.ImportFriend)
	}
	//group related routing group
	groupRouterGroup := r.Group("/group")
	{
		groupRouterGroup.POST("/create_group", group.CreateGroup)
		groupRouterGroup.POST("/set_group_info", group.SetGroupInfo)
		groupRouterGroup.POST("join_group", group.JoinGroup)
		groupRouterGroup.POST("/quit_group", group.QuitGroup)
		groupRouterGroup.POST("/group_application_response", group.ApplicationGroupResponse)
		groupRouterGroup.POST("/transfer_group", group.TransferGroupOwner)
		groupRouterGroup.POST("/get_group_applicationList", group.GetGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", group.GetGroupsInfo)
		groupRouterGroup.POST("/kick_group", group.KickGroupMember)
		groupRouterGroup.POST("/get_group_member_list", group.GetGroupMemberList)
		groupRouterGroup.POST("/get_group_all_member_list", group.GetGroupAllMember)
		groupRouterGroup.POST("/get_group_members_info", group.GetGroupMembersInfo)
		groupRouterGroup.POST("/invite_user_to_group", group.InviteUserToGroup)
		groupRouterGroup.POST("/get_joined_group_list", group.GetJoinedGroupList)
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister)
		authRouterGroup.POST("/user_token", apiAuth.UserToken)
	}
	//Third service
	thirdGroup := r.Group("/third")
	{
		thirdGroup.POST("/tencent_cloud_storage_credential", apiThird.TencentCloudStorageCredential)
	}
	//Message
	chatGroup := r.Group("/msg")
	{
		chatGroup.POST("/newest_seq", apiChat.GetSeq)
		chatGroup.POST("/pull_msg", apiChat.PullMsg)
		chatGroup.POST("/send_msg", apiChat.SendMsg)
		chatGroup.POST("/pull_msg_by_seq", apiChat.PullMsgBySeqList)
	}
	//Manager
	managementGroup := r.Group("/manager")
	{
		managementGroup.POST("/delete_user", manage.DeleteUser)
		managementGroup.POST("/send_msg", manage.ManagementSendMsg)
		managementGroup.POST("/get_all_users_uid", manage.GetAllUsersUid)
		managementGroup.POST("/account_check", manage.AccountCheck)
		managementGroup.POST("/get_users_online_status", manage.GetUsersOnlineStatus)
	}
	//Conversation
	conversationGroup := r.Group("/conversation")
	{
		conversationGroup.POST("/set_receive_message_opt", conversation.SetReceiveMessageOpt)
		conversationGroup.POST("/get_receive_message_opt", conversation.GetReceiveMessageOpt)
		conversationGroup.POST("/get_all_conversation_message_opt", conversation.GetAllConversationMessageOpt)
	}

	log.NewPrivateLog("api")
	ginPort := flag.Int("port", 10000, "get ginServerPort from cmd,default 10000 as port")
	flag.Parse()
	r.Run(":" + strconv.Itoa(*ginPort))
}
