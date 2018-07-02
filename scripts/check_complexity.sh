#!/bin/sh

print_help () {
        echo "A simple shell script to verify code complexity"
        echo ""
        echo "Usage: sh check_complexity.sh [-t,--threshold] <threshold> [-f,--files] <file(s)>"
        echo ""
        echo "Options:"
	echo "	-t|--threshold		The cyclomatic compelxity threshold you wish to test against."
	echo "	-f|--files		The file or set of files you're testing."
	echo ""
}

while [ $# -gt 0 ]
do
key=$1

case $key in
        -f|--files)
        files=$2
        shift
        ;;

        -t|--threshold)
        threshold=$2
        shift
        ;;


        -h|--help)
        print_help
        exit 0
        ;;

esac
shift
done

if [ -z "${threshold}" ]; then
        echo "-t|--threshold is a required input parameters"
        print_help
        exit 1
elif [ -z "${files}" ]; then
	echo "-f|--files is a required input parameters"
	print_help
	exit 1
fi

old_fs=$IFS
IFS='
'

cyclo_res=$(gocyclo ${files})

return_status=0

for res in ${cyclo_res} 
do
	complexity=$(echo ${res} | awk '{print $1}')
	if [ ${complexity} -ge ${threshold} ]; then
		echo "Function too complex: ${res}"
		return_status=1
	fi
done

IFS=${old_fs}

exit ${return_status}
