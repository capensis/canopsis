#!/bin/bash
set -e

branch1=$1
branch2=$2

if [ -z $branch1 ] || [ -z $branch2 ]; then
	echo "Usage: detect.sh branch1 branch2"
	exit 0
fi 

echo " :: Analyse diff between $branch1 and $branch2"
package=""
for tbl in $(git diff --name-only $branch1...$branch2 sources/); do 
	_package=$(echo $tbl | cut -d '/' -f2)

	if [ "$package" != "$_package" ]; then
		if [ -d "sources/$_package" ]; then
			if ! [[ $_package =~ (extra|externals|build.d|packages) ]]; then
				package="$_package"
				packages_list="$packages_list $_package"
			fi
		fi
	fi
done

package=""
echo " :: Analyse externals folder"
for tbl in $(git diff --name-only $branch1...$branch2 sources/externals); do 
	tbl=$(echo $tbl | sed 's|sources/externals/||')
	_package=$(echo $tbl | cut -d '-' -f1)

	if [ "$package" != "$_package" ]; then
		package="$_package"
		packages_list="$packages_list $_package"
	fi
done

package=""
echo " :: Analyse build.d folder"
for tbl in $(git diff --name-only $branch1...$branch2 sources/build.d); do 
	tbl=$(echo $tbl | sed 's|sources/build.d/||')
	_package=$(echo $tbl | sed 's|.install$||g')
	_package=$(echo $_package | sed 's|^[0-9][0-9\]_||g')

	if [ "$package" != "$_package" ]; then
		package="$_package"
		packages_list="$packages_list $_package"
	fi
done

package=""
echo " :: Analyse packages folder"
for tbl in $(git diff --name-only $branch1...$branch2 sources/packages); do 
	_package=$(echo $tbl | sed 's|sources/packages/||')

	if [ "$package" != "$_package" ]; then
		if [ -d "sources/$_package" ]; then
		package="$_package"
		packages_list="$packages_list $_package"
		fi
	fi
done

echo " :: Clean list"
packages_list=$(echo $packages_list | sed 's| |\n|g' | sort -u)

echo " :: Packages need to be recompiled"
for package in $packages_list; do
	echo "    - $package"
done

echo " :: Take a look at extra folder"
for tbl in $(git diff --name-only $branch1...$branch2 sources/extra); do 
	echo "    + $tbl " | sed 's|sources/extra/||'
done

