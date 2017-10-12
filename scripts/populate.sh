#!/bin/bash
# usage:
#	chmod +x populate.sh
#	./populate.sh [DATA] [IP] [PORT]
#
# DATA needs to be a plain text file with one entry value per line i.e.
# 1. headword
# 2. wordtype
# 3. definition
# 4. dictionary language (for the headword)
# 5. target language (for the definition)
#
index=0
count=0
word=""
wtyp=""
defn=""
dict=""
lang=""
ip=$2
port=$3
while IFS='' read -r line || [[ -n "$line" ]]; do
	case $index in
		0) word=$line
		((count++))
		;;
		1) wtyp=$line
		;;
		2) defn=$line
		;;
		3) dict=$line
		;;
		4) lang=$line
		body=$(cat <<EOF
{
  "headword": "$word",
  "wordtype": "$wtyp",
  "definition": "$defn",
  "hw_lang": "$dict",
  "def_lang": "$line"
}
EOF
)
		echo "REQUEST #$count"
		echo "$body"
		curl -X POST -H "Content-Type: application/json" -d "$body" "http://$ip:$port/entry"
		index=-1
		;;
	esac
	((index++))
done < "$1"
echo "Processed $count dictionary entry(s)."
