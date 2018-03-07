echo $PWD | grep src$ && { 
    echo "Cannot run clean in src directory!"
    exit 1
}

for i in `ls | egrep -v ^src`
do
    echo rm -rf ${i}
done
