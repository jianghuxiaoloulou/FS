package models

import (
	"WowjoyProject/FileServer/pkg/loggin"
	"WowjoyProject/FileServer/pkg/setting"
	"os"
	"path/filepath"
	"strings"
)

type ActionType int

const (
	UPLOAD ActionType = iota
	DOWNLOAD
	DELETE
)

type ObjectData struct {
	InstanceKey  int64
	Key          string
	Type         ActionType
	SyncStrategy string
	Path         string
}

// 对象存储
// web单文件上传
func UploadFile(file string) bool {
	return true
}

// web检查号上传
func UploadNumbers(number string) bool {
	return true
}

// web单文件下载
func DownFile(file string) string {
	return ""
}

// web检查号下载
func DownNumbers(number string) string {
	return ""
}

// web单文件删除
func DeleteFile(file string) bool {
	return true
}

// web检查号删除
func DeleteNumbers(number string) bool {
	return true
}

// 自动上传文件
// 获取需要上传的数据
func AutoUploadData(dataChan chan ObjectData) {
	sql := `select im.instance_key,im.img_file_name, ins.file_name,stu.ip,stu.s_virtual_dir
	from  image im 
	inner join instance ins on im.instance_key = ins.instance_key
	inner join study_location stu on ins.location_code = stu.n_station_code
	where im.dcm_file_upload_status = 1 order by im.instance_key ASC limit ?;`
	rows, err := db.Query(sql, setting.MaxTasks)
	if err != nil {
		loggin.Fatal(err)
		return
	} else {
		for rows.Next() {
			var instance_key int64
			var imgfile, dcmfile, ip, virpath string
			_ = rows.Scan(&instance_key, &imgfile, &dcmfile, &ip, &virpath)
			if imgfile != "" {
				filefullpath := fileFullPath(imgfile, ip, virpath)
				loggin.Info("需要上传的文件名：", filefullpath)
				data := ObjectData{
					InstanceKey:  instance_key,
					Key:          imgfile,
					Type:         UPLOAD,
					SyncStrategy: setting.OBJECT_Sync,
					Path:         filefullpath,
				}
				dataChan <- data
			}
			if dcmfile != "" {
				filefullpath := fileFullPath(dcmfile, ip, virpath)
				loggin.Info("需要上传的文件名：", filefullpath)
				data := ObjectData{
					InstanceKey:  instance_key,
					Key:          dcmfile,
					Type:         UPLOAD,
					SyncStrategy: setting.OBJECT_Sync,
					Path:         filefullpath,
				}
				dataChan <- data
			}
		}
		rows.Close()
	}
}

func fileFullPath(file, ip, virpath string) (path string) {
	if file == "" || ip == "" || virpath == "" {
		path = ""
	} else {
		path += "\\\\"
		path += ip
		path += "\\"
		path += virpath
		path += "\\"
		path += file
	}
	return
}

func getFileSize(filename string) int64 {
	var result int64
	if Exist(filename) {
		filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
			result = f.Size()
			return nil
		})
	}
	return result
}

func FullFilePath(file, path string) (fullpath string) {
	if file == "" || path == "" {
		fullpath = ""
	} else {
		fullpath += path
		fullpath += file
	}
	return
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 上传数据后更新数据库
func UpdateAutoUplaode(key int64, file string, value bool) {
	// value true 上传成功
	// value false 上传失败
	if value {
		if strings.Contains(file, ".dcm") {
			sql := `update image im set im.dcm_file_upload_status=0,im.dcm_file_name_remote=? where im.instance_key=?`
			db.Exec(sql, file, key)
		} else {
			sql := `update image im set im.img_file_upload_status=0,im.img_file_name_remote=? where im.instance_key=?`
			db.Exec(sql, file, key)
		}
	} else {
		if strings.Contains(file, ".dcm") {
			sql := `update image im set im.dcm_file_upload_status=10001 where im.instance_key=? and im.dcm_file_upload_status!=0`
			db.Exec(sql, key)
		} else {
			sql := `update image im set im.img_file_upload_status=10001 where im.instance_key=? and im.img_file_upload_status!=0`
			db.Exec(sql, key)
		}
	}
}

func TestAutoDownData(dataChan chan ObjectData) {
	path := setting.OBJECT_PATH + "1a3a78e6616ebcbeafd5f577432d1264.dcm"
	data := ObjectData{
		InstanceKey:  1111,
		Key:          path,
		Type:         DOWNLOAD,
		SyncStrategy: setting.OBJECT_Sync,
		Path:         "W:\\image\\1a3a78e6616ebcbeafd5f577432d1264.dcm",
	}
	dataChan <- data
}

// 自动下载任务:
func AutoDownData(dataChan chan ObjectData) {
	loggin.Info("开始获取下载数据......")
	sql := `select im.instance_key,im.img_file_name,im.img_file_name_remote,im.dcm_file_name_remote,ins.file_name
	from  image im 
	inner join instance ins on im.instance_key = ins.instance_key
	where ins.FileExist = -1 order by ins.instance_key limit ?;`
	rows, err := db.Query(sql, setting.MaxTasks)
	if err != nil {
		loggin.Fatal(err)
		return
	} else {
		for rows.Next() {
			data := DownData{}
			err = rows.Scan(&data.instance_key, &data.jpgfile, &data.jpgremote, &data.dcmremote, &data.dcmfile)
			if err != nil {
				loggin.Error(err)
			}
			if data.dcmfile.Valid && data.dcmfile.String != "" {
				fullpath := FullFilePath(data.dcmfile.String, setting.DestRoot)
				remotepath := setting.OBJECT_PATH + data.dcmremote.String
				size := getFileSize(fullpath) / 1024
				// 判断已经下载文件大小2KB
				if size < 2 {
					loggin.Info("需要下载的文件名：", data.dcmfile.String)
					data := ObjectData{
						InstanceKey:  data.instance_key.Int64,
						Key:          remotepath,
						Type:         DOWNLOAD,
						SyncStrategy: setting.OBJECT_Sync,
						Path:         fullpath,
					}
					dataChan <- data
				} else {
					loggin.Info("文件已经存在，直接更新:", fullpath)
					UpdateAutoDown(data.instance_key.Int64, remotepath, true)
				}

			}
			if data.jpgfile.Valid && data.jpgfile.String != "" {
				fullpath := FullFilePath(data.jpgfile.String, setting.DestRoot)
				size := getFileSize(fullpath) / 1024
				// 判断已经下载文件大小2KB
				if size < 2 {
					loggin.Info("需要下载的文件名：", data.jpgfile.String)
					remotepath := setting.OBJECT_PATH + data.jpgremote.String
					data := ObjectData{
						InstanceKey:  data.instance_key.Int64,
						Key:          remotepath,
						Type:         DOWNLOAD,
						SyncStrategy: setting.OBJECT_Sync,
						Path:         fullpath,
					}
					dataChan <- data
				}
			}
		}
		rows.Close()
	}
}

// 数据下载成功更新数据库
func UpdateAutoDown(key int64, file string, value bool) {
	// value true 上传成功
	// value false 上传失败
	code := setting.DestCode
	if value {
		if strings.Contains(file, ".dcm") {
			sql := `update instance ins set ins.FileExist = 2,ins.location_code=? where ins.instance_key=?`
			db.Exec(sql, code, key)
		}
	} else {
		if strings.Contains(file, ".dcm") {
			sql := `update instance ins set ins.FileExist = -2 where ins.instance_key=?`
			db.Exec(sql, key)
		}
	}
}
