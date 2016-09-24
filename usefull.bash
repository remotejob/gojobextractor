scp 104.236.237.125:gojobextractor/mytags.csv .
scp 104.236.237.125:gojobextractor/notmytags.csv .





GOBIN=$PWD/bin go install


 GOBIN=$PWD/bin go install apply_for_job/apply_headless/apply_headless.go
GOBIN=$PWD/bin go install sendemailtoemploer/sendemailtoemploer.go 

stop mongodb

mongodump --dbpath  /var/lib/mongodb
mongorestore --dbpath /var/lib/mongodb dump

start mongodb



rm -rf /tmp/cv_employers
mongodump --db=cv_employers --out=/tmp
scp -rp /tmp/cv_employers juno@104.236.240.215:/tmp   #remotejob
scp -rp /tmp/cv_employers juno@104.236.237.125:/tmp  #cv


mongorestore --db=cv_employers /tmp/cv_employers




export PATH=$PATH:/home/juno/selenium
GOPATH=$GOPATH:~/neonworkspace/gojobextractor go test -v

$GOPATH/bin/ginkgo bootstrap --agouti

GOPATH=$GOPATH:/home/juno/git/jobprotractor/gojobextractor /home/juno/workspace/gocode/bin/ginkgo generate --agouti user_login


db.getCollection('employers').find({"created_at":{ 
    $lt: new Date(), 
    $gte: new Date(new Date().setDate(new Date().getDate()-1))
  }}   )

db.getCollection('employers').find({"created_at" : { $gte : new ISODate("2016-05-08T00:00:00Z") },"applied":false})
db.getCollection('employers').find({"created_at" : { $gte : new ISODate("2016-05-08T00:00:00Z") },"externallink":"","applied":false})
db.getCollection('employers').find({"created_at" : { $gte : new ISODate("2016-05-08T00:00:00Z") },"externallink":/mailto/})
db.getCollection('employers').find({"created_at" : { $gte : new ISODate("2016-05-08T00:00:00Z") },"location":/Finland/})
db.getCollection('employers').find({"externallink":{$ne:""},"location":/Finland/})
db.getCollection('employers').find({"externallink":{$ne:""},"email":{$ne:""},"applied":false,"location":/Finland/})
db.getCollection('employers').find({"externallink":{$ne:""},"email":{$ne:""},"applied":false})
db.getCollection('employers').find({"externallink":"mailto:jobs@nitor.fi?subject=Full%20Stack%20Developer%20(via%20Stack%20Overflow%20Careers)&body=%0d%0a--%0d%0aFound%20via%20Stack%20Overflow%20Careers%0d%0a"}

db.getCollection('employers').find({"location":/Finland/,"externallink":{"$ne": ""},"applied":false})
db.getCollection('employers').find({"externallink":/humany/})


fresh (for server)


/usr/local/bin/phantomjs  --webdriver-selenium-grid-hub=http://192.168.1.4:4444/grid/register/
Xvfb :10 -ac



export DISPLAY=:11
Xvfb :11 -screen 0 1024x768x16 &

java -jar selenium-server-standalone-2.53.0.jar \
   -role node \
   -port 4441
   -hub http://localhost:4444/grid/register \
   -browser "browserName=firefox,version=19,maxInstances=5"&
