package model

var DemoVideos = []Video{
	{
		Id:            1,
		Path:          "https://www.w3schools.com/html/movie.mp4",
		CoverPath:     "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
}
