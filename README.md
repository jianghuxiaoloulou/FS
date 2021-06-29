# 文件上传下载服务
# 1.自动上传
# 1.1 通过image表中 dcm_file_upload_status字段判断待上传的状态来自动上传数据 ('0-上传成功，1-待上传，-1-默认错误码（海纳数据），-10001-上传失败, 2图像文件没有生成')
# <!-- CREATE TABLE dummy (
#	`img_file_upload_status` SMALLINT(6) NULL DEFAULT '2' COMMENT '0-上传成功，1-待上传，2 - 图像未生成 -1-默认错误码（海纳数据），-10001-上传失败'
# ) -->
# 上传文件后，修改文件状态，填写上传文件key
# 通过instance_key来关联image表和instance表，来获取dcm文件和jpg文件
# 增加医院id


# 2.http上传
# 2.1 新建表, 表中包含字段(key,uid_enc,accession_number,upload_status,down_status,remote_file_name,down_file_name,localtion_code,backet_id,create_time,update_time)


## create a new repository on the command line
# echo "# FS" >> README.md
# git init
# git add README.md
# git commit -m "first commit"
# git branch -M main
# git remote add origin https://github.com/jianghuxiaoloulou/FS.git
# git push -u origin main