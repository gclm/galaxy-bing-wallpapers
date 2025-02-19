package model

import (
	"fmt"
	"strings"
)

// Wallpaper 必应壁纸数据结构
type Wallpaper struct {
	ID            int    `bson:"id" json:"id"`                       // 唯一标识
	Title         string `bson:"title" json:"title"`                 // 图片标题
	Url           string `bson:"url" json:"url"`                     // 图片URL
	Datetime      string `bson:"datetime" json:"datetime"`           // 日期时间
	Copyright     string `bson:"copyright" json:"copyright"`         // 版权信息
	CopyrightLink string `bson:"copyrightlink" json:"copyrightlink"` // 版权链接
	Hsh           string `bson:"hsh" json:"hsh"`                     // 哈希值
	CreatedTime   string `bson:"created_time" json:"created_time"`   // 创建时间
	Mkt           string `bson:"mkt" json:"mkt"`                     // 市场代码 如：fr-FR
}

// WallpaperResponse API响应结构
type WallpaperResponse struct {
	Images []Wallpaper `json:"images"`
	Total  int64       `json:"total"`
}

// GenerateImageURL 生成指定尺寸的图片URL
func (w *Wallpaper) GenerateImageURL(width, height string) string {
	// 从原始URL中提取基础部分并替换尺寸
	baseURL := w.Url[:strings.LastIndex(w.Url, "_")+1]
	return baseURL + width + "x" + height + ".jpg"
}

// ToMap 将Wallpaper转换为bson.M
func (w *Wallpaper) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            w.ID,
		"title":         w.Title,
		"url":           w.Url,
		"datetime":      w.Datetime,
		"copyright":     w.Copyright,
		"copyrightlink": w.CopyrightLink,
		"hsh":           w.Hsh,
		"created_time":  w.CreatedTime,
		"mkt":           w.Mkt,
	}
}

// Validate 验证壁纸数据
func (w *Wallpaper) Validate() error {
	if w.Title == "" {
		return fmt.Errorf("title is required")
	}
	if w.Url == "" {
		return fmt.Errorf("url is required")
	}
	if w.Datetime == "" {
		return fmt.Errorf("datetime is required")
	}
	if w.Mkt == "" {
		return fmt.Errorf("mkt is required")
	}
	return nil
}
