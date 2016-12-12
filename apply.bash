#! /bin/bash

if [ -f "./cant_apply.txt" ]; then ./bin/cant_apply;fi

# cp config.gcfg.alastonkuvia.fi config.gcfg && 
# ./bin/apply_headless && cp config.gcfg.alastonkuvia.info config.gcfg && 
# ./bin/apply_headless


cp config.gcfg.astrologi.fi config.gcfg && 
./bin/apply_headless && cp config.gcfg.antaa.fi config.gcfg && 
./bin/apply_headless
