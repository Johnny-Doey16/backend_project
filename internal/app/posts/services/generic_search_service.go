package services

import (
	"encoding/json"
	"time"

	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserD struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	Name     string `json:"name"`
	// LastName       string    `json:"last_name"`
	HeaderImageUrl string    `json:"header_image_url"`
	ImageUrl       string    `json:"image_url"`
	About          string    `json:"about"`
	IsVerified     bool      `json:"is_verified"`
	IsFollowing    bool      `json:"is_following"`
	CreatedAt      time.Time `json:"created_at"`
}

type ChurchD struct {
	Id             string    `json:"id"`
	ChurchUsername string    `json:"church_username"`
	Name           string    `json:"name"`
	MembersCount   int64     `json:"members_count"`
	HeaderImageUrl string    `json:"header_image_url"`
	UserType       string    `json:"user_type"`
	ImageUrl       string    `json:"image_url"`
	About          string    `json:"about"`
	IsMember       bool      `json:"is_member"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
	ChurchId       int64     `json:"church_id"`
	IsFollowing    bool      `json:"is_following"`
	City           string    `json:"city"`
	State          string    `json:"state"`
}

type PostD struct {
	PostUserId string `json:"post_user_id"`
	Id         string `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	// LastName    string    `json:"last_name"`
	Content     string    `json:"content"`
	IsVerified  bool      `json:"is_verified"`
	UserType    string    `json:"user_type"`
	About       string    `json:"about"`
	ImageUrl    string    `json:"image_url"`
	PostImages  []string  `json:"post_images"`
	Comments    int64     `json:"comments"`
	CreatedAt   time.Time `json:"created_at"`
	TotalImages int64     `json:"total_images"`
	Views       int64     `json:"views"`
	Likes       int64     `json:"likes"`
	Repost      int64     `json:"repost"`
	PostLiked   bool      `json:"post_liked"`
}

func UnmarshalUserResult(data []byte) ([]*pb.UserData, error) {
	var users []UserD
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	var pbUsers []*pb.UserData
	for _, user := range users {
		pbUsers = append(pbUsers, convertUserDtoToPb(user))
	}

	return pbUsers, nil
}

func UnmarshalChurchResult(data []byte) ([]*pb.ChurchData, error) {
	var churches []ChurchD
	if err := json.Unmarshal(data, &churches); err != nil {
		return nil, err
	}

	var pbChurches []*pb.ChurchData
	for _, church := range churches {
		pbChurches = append(pbChurches, convertChurchDtoToPb(church))
	}

	return pbChurches, nil
}

func UnmarshalPostResult(data []byte) ([]*pb.PostData, error) {
	var posts []PostD
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}

	var pbPosts []*pb.PostData
	for _, post := range posts {
		pbPosts = append(pbPosts, convertPostDtoToPb(post))
	}

	return pbPosts, nil
}

func convertUserDtoToPb(user UserD) *pb.UserData {
	return &pb.UserData{
		Id:             user.Id,
		Username:       user.Username,
		UserType:       user.UserType,
		Name:           user.Name,
		ImageUrl:       user.ImageUrl,
		IsVerified:     user.IsVerified,
		HeaderImageUrl: user.HeaderImageUrl,
		About:          user.About,
		IsFollowing:    user.IsFollowing,
		CreatedAt:      timestamppb.New(user.CreatedAt),
		// LastName:       user.LastName,
	}
}

func convertChurchDtoToPb(church ChurchD) *pb.ChurchData {
	return &pb.ChurchData{
		UserData: &pb.UserData{
			Id:             church.Id,
			ImageUrl:       church.ImageUrl,
			IsVerified:     church.IsVerified,
			HeaderImageUrl: church.HeaderImageUrl,
			CreatedAt:      timestamppb.New(church.CreatedAt),
			About:          church.About,
			Username:       church.ChurchUsername,
			UserType:       church.UserType,
			IsFollowing:    church.IsFollowing,
		},
		Name:         church.Name,
		MembersCount: church.MembersCount,
		IsMember:     church.IsMember,
		ChurchId:     church.ChurchId,
		Location: &pb.LocationData{
			City:  church.City,
			State: church.State,
		},
	}
}

func convertPostDtoToPb(post PostD) *pb.PostData {
	// log.Println("Post ID", post.Id, "UID", post.PostUserId, "Like", post.PostLiked, "Created at", post.CreatedAt)
	return &pb.PostData{
		UserData: &pb.UserData{
			Id:         post.PostUserId,
			ImageUrl:   post.ImageUrl,
			IsVerified: post.IsVerified,
			Username:   post.Username,
			Name:       post.Name,
			About:      post.About,
			UserType:   post.UserType,
			// LastName:   post.LastName,
			// CreatedAt:  timestamppb.New(post.CreatedAt),
		},
		Metrics: &pb.PostMetricsData{
			Comments:  post.Comments,
			Likes:     post.Likes,
			Views:     post.Views,
			Reposts:   post.Repost,
			PostLiked: post.PostLiked,
		},
		Id:          post.Id,
		Content:     post.Content,
		PostImages:  post.PostImages,
		TotalImages: int64(post.TotalImages),
		CreatedAt:   timestamppb.New(post.CreatedAt),
	}
}
