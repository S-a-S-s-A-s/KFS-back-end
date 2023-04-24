package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Id       string
}

type Page struct {
	PageNum  int
	PageSize int
	Keyword  string
	Desc     bool
}

func Info(c *gin.Context) {
	// name := c.GetHeader("userName")
	token := c.GetHeader("token")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"roles":  "[admin]",
			"name":   token_username[token].UserName,
			"avatar": "https://oss.aliyuncs.com/aliyun_id_photo_bucket/default_handsome.jpg",
		},
	})
}

func Login(c *gin.Context) {
	//fmt.Println(c)
	var user User
	c.BindJSON(&user)
	//fmt.Println(user)
	token, id := TokensTem(user.UserName, user.PassWord)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "权限不足",
			"data": gin.H{
				"name":  user.UserName,
				"token": token,
			},
		})
	} else {
		user.Id = id
		token_username[token] = user
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"name":  user.UserName,
				"token": token,
			},
		})
	}
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "success",
	})
}

// 分页查询
func PageAndLimit(c *gin.Context) {
	var p Page
	p.PageNum, _ = strconv.Atoi(c.Param("page"))
	p.PageSize, _ = strconv.Atoi(c.Param("limit"))

	name := c.Query("name")
	project := c.Query("projectName")
	projectId := c.Query("projectId")
	//fmt.Println(projectId, project)
	token := c.GetHeader("token")
	//fmt.Println(token)
	role := GetRole(projectId, token)
	//fmt.Println(p)
	if p.PageNum <= 0 {
		p.PageNum = 1
	}
	var files []File
	var query string
	query = GetQuery(role, name, project, token_username[token].UserName)
	//fmt.Println(query)
	var count int64
	DB.Model(&File{}).Where(query).Count(&count)
	if err := DB.Model(&File{}).Where(query).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": gin.H{"records": files,
			"total": count,
		},
	})
}

// 下载文件
func DownloadFile(c *gin.Context) {
	token := c.PostForm("token")
	name := c.Param("Name")
	project := c.Param("Project")
	user := token_username[token]
	c.File("./file/" + project + "/" + name)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "下载成功",
	})
	var file = File{}
	DB.Where("name = ? AND project = ?", name, project).First(&file)
	DB.Create(&FileInfo{
		Operation: "download",
		Operator:  user.UserName,
		FileId:    file.ID,
		Status:    true,
		CreatTime: time.Now().Format("2006-01-02 15:04:05"),
	})
}

// 上传文件
func UploadFile(c *gin.Context) {
	token := c.PostForm("token")
	// fmt.Println(token)
	// 获取文件头
	file, err := c.FormFile("file")
	projectName := c.PostForm("projectName")
	isPublic := c.PostForm("isPublic")
	var ispublic bool
	if isPublic == "是" {
		ispublic = true
	}
	user := token_username[token]
	if err != nil {
		fmt.Println(err.Error())
	}
	fileName := file.Filename
	err = c.SaveUploadedFile(file, "./file/"+projectName+"/"+fileName)
	var fileTem = File{}

	DB.Create(&File{
		Name:      fileName,
		Size:      formatFileSize(file.Size),
		Project:   projectName,
		UserName:  user.UserName,
		IsPublic:  ispublic,
		CreatTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	DB.Where("name = ? AND project = ?", fileName, projectName).First(&fileTem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "上传失败",
		})
		DB.Create(&FileInfo{
			Operation: "upload",
			Operator:  user.UserName,
			FileId:    fileTem.ID,
			Status:    false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功",
	})

	DB.Create(&FileInfo{
		Operation: "upload",
		Operator:  user.UserName,
		FileId:    fileTem.ID,
		Status:    true,
		CreatTime: time.Now().Format("2006-01-02 15:04:05"),
	})
}

// 得到该用户所在的所有项目
func GetProject(c *gin.Context) {
	token := c.GetHeader("token")
	pros := GetProjects(token)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": pros,
	})
}

// 获得文件大小
func GetFileSize(c *gin.Context) {
	name := c.Param("Name")
	project := c.Param("Project")
	fi, err := os.Stat("./file/" + project + "/" + name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"size": fi.Size(),
			},
		})
	}
}

// 大文件下载
func DownloadBigFile(c *gin.Context) {
	projectName := c.PostForm("projectName")
	fileName := c.PostForm("fileName")
	chunkSize, _ := strconv.ParseInt(c.PostForm("chunkSize"), 10, 64)
	index, _ := strconv.ParseInt(c.PostForm("index"), 10, 64)
	chunkTotal, _ := strconv.ParseInt(c.PostForm("chunkTotal"), 10, 64)
	//file, _ := os.ReadFile("./file/" + projectName + "/" + fileName)

	//chunk := file[offset : offset+chunkSize]
	if index == 1 {
		token := c.PostForm("token")
		user := token_username[token]
		var fileTem = File{}
		DB.Where("name = ? AND project = ?", fileName, projectName).First(&fileTem)
		DB.Create(&FileInfo{
			Operation: "download",
			Operator:  user.UserName,
			FileId:    fileTem.ID,
			Status:    false,
			CreatTime: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	fi, _ := os.Open("./file/" + projectName + "/" + fileName)
	defer fi.Close()
	offset := chunkSize * (index - 1)
	fileInfo, _ := fi.Stat()
	if index == chunkTotal {
		offset = fileInfo.Size() - chunkSize
	}
	fi.Seek(offset, 0)
	r := bufio.NewReader(fi)
	chunk := make([]byte, chunkSize)
	r.Read(chunk)
	c.Header("filename", fileName)
	c.Header("Content-Length", strconv.Itoa(len(chunk)))
	c.Header("Content-Disposition", "attachment;filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(200, "application/octet-stream", chunk)
	if index == chunkTotal {
		token := c.PostForm("token")
		user := token_username[token]
		var fileTem = File{}
		DB.Where("name = ? AND project = ?", fileName, projectName).First(&fileTem)
		var fileInfoTem FileInfo
		DB.Where("file_id = ? AND operator = ?", fileTem.ID, user.UserName).Order("id desc").First(&fileInfoTem)
		DB.Model(&FileInfo{ID: fileInfoTem.ID}).Update("status", true)
	}
}

func DownloadAuth(c *gin.Context) {

}

func UploadAuth(c *gin.Context) {
	token := c.GetHeader("token")
	projectId := c.Query("projectId")
	role := GetRole(projectId, token)
	if role == "manager" || role == "vip" {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "通过",
		})
		return
	}
	c.JSON(203, gin.H{
		"code": 203,
		"msg":  "你的权限不足",
	})
}
