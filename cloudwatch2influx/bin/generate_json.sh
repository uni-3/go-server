#!/bin/bash

# 引数の数は1個
if [ $# -ne 1 ]; then
  echo "arguments are $# items" 1>&2
  echo "1 arguments(how many generate json row) required to execute" 1>&2
  exit 1
fi

# 引数で指定した行数だけ生成
for l in `seq $1`; do
  echo "{\"idle\": `awk 'BEGIN{ srand('"$RANDOM"'); printf("%2.1f", rand()*100) }'`, \
\"system\": `awk 'BEGIN{ srand('"$RANDOM"'); printf("%2.1f", rand()*100) }'`, \
\"user\": `awk 'BEGIN{ srand('"$RANDOM"'); printf("%2.1f", rand()*100) }'`\
}"
done

# 同じ数が連続して出てくるのでやめた
#yes "echo {\"cpu\": $((RANDOM%+101)), }" | head -10 | bash
