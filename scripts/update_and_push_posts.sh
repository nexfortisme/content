./generate_index_json.sh

current_datetime="$(date '+%Y-%m-%d %H:%M:%S')" 
msg="Post Update at $current_datetime"

git add -A
git c -m "$msg"
git pu

