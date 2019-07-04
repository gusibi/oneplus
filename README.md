## code2png
markdown code 替换为图片

### 使用方法

```
git clone git@github.com:gusibi/oneplus.git

cd oneplus/code2png

python plus.py -m [markdown_path] -n [outfile_path]

```

### 示例 

这是转换前：

https://github.com/gusibi/oneplus/blob/master/325.md

这是转换后：

https://github.com/gusibi/oneplus/blob/master/new_325.md

## html-server

该目录为 code html 页面渲染服务

## mobile-attribution

手机号归属地查询，服务部署在腾讯云serverless 服务，使用ApiGateway 触发

### restful API

serverless 提供服务

### 数据更新

定时任务触发
