#!/bin/bash

# This script convert tests from the format of infos to the pesho.org format.
# Pass as first argument a .zip argive downloaded from infos. The tool will
# unzip, find all folders with tests, convert the tests and zip them again.
# As a result one .zip per test folder will be created.

base_dir=${PWD}	

if [ -z "${1}" ]
then
	echo "Plaese provide an infos archive as first argument."
	exit 1
fi

echo "Working with ${1}"
unzip "${1}"

for dir in `find . -type d` ; do
     	echo "${dir}"
     	all_files=`find "${dir}" -type f | wc -l`
     	test_files=`find "${dir}" -name "*.sol" -o -name "*.in" | wc -l`
	echo "all files ${all_files}"
	echo "test files ${test_files}"
	if [[ "${all_files}" == "${test_files}" && "${test_files}" > 0 ]]
	then
		counter=1
		for f in `find "${dir}" -name "*.in"`; do 
			echo "$f"
			mv  "$f" "${dir}/input$counter"
			counter=$((counter+1))
		done

		counter=1
		for f in `find "${dir}" -name "*.sol"`; do
			echo "$f"
			mv  "$f" "${dir}/output$counter"
			counter=$((counter+1))
		done
		
		echo "Zipping ${dir}"
		cd "${dir}"
		zip "tests.zip" `find . -name "input*" -o -name "output*"`
		cd "${base_dir}"
	fi
done



