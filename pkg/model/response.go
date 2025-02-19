package model

// WallpaperList 壁纸列表响应结构
type WallpaperList struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Total int         `json:"total"`
	Data  []Wallpaper `json:"data"`
}

// ImageResponse 图片信息响应结构
type ImageResponse struct {
	Url      string `json:"url"`      // 图片URL
	Title    string `json:"title"`    // 图片标题
	Datetime string `json:"datetime"` // 日期时间
}

// 添加统一的API响应结构
type ApiResponse struct {
	Code    int         `json:"code"`            // 状态码
	Message string      `json:"message"`         // 响应信息
	Data    interface{} `json:"data,omitempty"`  // 响应数据
	Total   int64       `json:"total,omitempty"` // 总数（列表接口使用）
}
