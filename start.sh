# !/bin/bash
# FactoCord start script
# By Allen Lydiard (FM1337)

rm /home/fact/fact-g/saves/*.tmp.*
screen -dmS FactoCord-g ./FactoCord
echo "Hey if you're running this manually and not as a service, you can go ahead and just ctrl+c this script and the server will continue to run!"
while screen -list | grep -q FactoCord-g
do
    sleep 30
done
ExitCode=`cat .exit`
exit $ExitCode
