pipeline {
  agent {
    label 'k8s-cf-go'
  }
  environment {
    TMP_DIR = '.TMP'
    BUILD_DIR = '.BUILD'
    CF_ORG = 'Audi Beta Space'
    CF_SPACE_PROD = 'PROD'
    CF_SPACE_DEV = 'DEV'
    MANIFEST_DEV = 'manifest.dev.yml'
    MANIFEST_PROD = 'manifest.prod.yml'
  }
  stages {
    stage('Setup') {
      steps {
        sh 'env'
        sh "mkdir ~/.ssh"
        sh "echo 'Host *\n\tStrictHostKeyChecking no\n' > ~/.ssh/config"
        sh 'go version'
        sh 'go2xunit --version'
        sh 'cf -v'
        sh 'echo $TMP_DIR'
        sh 'rm -rf $TMP_DIR && mkdir -p $TMP_DIR'
        sh 'rm -rf $BUILD_DIR && mkdir -p $BUILD_DIR'
      }
    }
    stage('Dependencies') {
      steps {
        sshagent(['e13d693f-202e-4f42-b708-ad46f034eb27']) {
          sh 'go mod vendor'
          }  
      }
    }
    stage('Lint') {
      steps {
        sh 'go vet ./...'
      }
    }    
    stage('Test') {
      steps {
        sh '2>&1 gocov test -v `go list ./... | grep -v "generated"` > $TMP_DIR/coverage.json | go2xunit -output $TMP_DIR/tests.xml'
        junit '**/tests.xml'
      }
    }
    stage('Health Check') {
      steps {
        sh 'cat $TMP_DIR/coverage.json | gocov-xml > $TMP_DIR/coverage.xml'
        sh 'gocov-html $TMP_DIR/coverage.json > $TMP_DIR/coverage.html'
        publishHTML target: [
            allowMissing: false,
            alwaysLinkToLastBuild: false,
            keepAll: true,
            reportDir: '.TMP',
            reportFiles: 'coverage.html',
            reportName: 'GoCove Report'
          ]
        cobertura autoUpdateHealth: false,
  				autoUpdateStability: false,
  				coberturaReportFile: '**/coverage.xml',
  				lineCoverageTargets: '80, 60, 0',
          conditionalCoverageTargets: '80, 60, 0',
  				methodCoverageTargets: '80, 60, 0',
  				maxNumberOfBuilds: 0,
  				onlyStable: true,
  				sourceEncoding: 'ASCII',
  				zoomCoverageChart: false,
  				failNoReports: true,
  				failUnhealthy: true,
  				failUnstable: true
      }
    }
     stage('Build') {
      steps {
  	    sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags \'-extldflags "-static"\' -o $BUILD_DIR/bin'
      }
    }
    stage('Deploy Development') {
      when {
        branch 'develop'
      }
      steps {
        dir("${BUILD_DIR}") {
          sh "cp ../$MANIFEST_DEV ./"
          withCredentials([file(credentialsId: 'secretFileVWSTempCreds', variable: 'FILE')]) {
            sh 'source $FILE && curl --user $accessKey:$secretKey -d "grant_type=client_credentials" $tokenUrl > token.json'
            sh 'cat token.json | sed "s/{.*\\"access_token\\":\\"\\([^\\"]*\\).*}/\\1/g" > token.txt'
            sh 'source $FILE && curl -H "Accept: text/x-shellscript" -H "Authorization: Bearer $(cat token.txt)" $credsUrl > creds'
            sh 'rm ./token.*'
            sh 'source creds && rm creds && cf login -a "$CF_API_URL" -u "$CF_USER" -p "$CF_PASSWORD" -o "$CF_ORG" -s "$CF_SPACE_DEV"'
            sh "cf push -f $MANIFEST_DEV"
          }
        }
      }
    }
    stage('Deploy Production') {
      when {
        branch 'master'
      }
      steps {
        dir("${BUILD_DIR}") {
          sh "cp ../$MANIFEST_PROD ./"
          withCredentials([file(credentialsId: 'secretFileVWSTempCreds', variable: 'FILE')]) {
            sh 'source $FILE && curl --user $accessKey:$secretKey -d "grant_type=client_credentials" $tokenUrl > token.json'
            sh 'cat token.json | sed "s/{.*\\"access_token\\":\\"\\([^\\"]*\\).*}/\\1/g" > token.txt'
            sh 'source $FILE && curl -H "Accept: text/x-shellscript" -H "Authorization: Bearer $(cat token.txt)" $credsUrl > creds'
            sh 'rm ./token.*'
            sh 'source creds && rm creds && cf login -a "$CF_API_URL" -u "$CF_USER" -p "$CF_PASSWORD" -o "$CF_ORG" -s "$CF_SPACE_PROD"'
            sh "cf push -f $MANIFEST_PROD"
          }
        }
      }
    }
  }
}