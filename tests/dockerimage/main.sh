#!/bin/sh

echo "STDOUT: 'main.sh' started."

# Ensure output filepath
OUT=/var/data/out.`date +"%Y_%m%d_%H%M%S"`.txt
touch ${OUT}

# Self introduction
echo "[WHO AM I?] `whoami`" >> ${OUT}
echo "[UNAME] `uname -a`" >> ${OUT}

# Check Env variables
echo "DAAPTEST_FOO: ${DAAPTEST_FOO}" >> ${OUT}
echo "DAAPTEST_BAR: ${DAAPTEST_BAR}" >> ${OUT}

# Check mounted files
echo "Read from source file >>" >> ${OUT}
cat /var/data/source.txt >> ${OUT}
echo "<< End source file" >> ${OUT}

echo "STDOUT: 'main.sh' successfully finished."
