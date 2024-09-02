#!/bin/bash

# 提示用户输入版本号
read -p "请输入版本号（例如 v1.0.0）： " version

# 检查版本号是否为空
if [ -z "$version" ]; then
  echo "版本号不能为空"
  exit 1
fi

# 删除本地标签
if git tag -d "$version" 2>/dev/null; then
  echo "本地标签 $version 已删除"
else
  echo "本地标签 $version 不存在或无法删除"
fi

# 删除远程标签
if git push origin --delete "$version" 2>/dev/null; then
  echo "远程标签 $version 已删除"
else
  echo "远程标签 $version 不存在或无法删除"
fi

# 推送新的标签
if git tag "$version" && git push origin "$version"; then
  echo "标签 $version 已成功推送到远程仓库"
else
  echo "标签 $version 推送失败"
  exit 1
fi