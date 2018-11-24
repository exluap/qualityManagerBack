pipeline {
  agent {
    docker {
      image 'resin/raspberrypi3-golang'
      args ' -v /opt/passwd/passwd:/etc/passwd:ro -v /opt/passwd/group:/etc/group:ro -v /opt/passwd/shadow:/etc/shadow:ro'
    }

  }
  stages {
    stage('set worspace') {
      steps {
        sh 'export GOPATH=$WORKSPACE'
      }
    }
    stage('get depend') {
      steps {
        sh '''go get github.com/dgrijalva/jwt-go
go get github.com/rhysd/go-github-selfupdate/selfupdate
go get github.com/rs/cors
go get github.com/mattn/go-sqlite3'''
        sh 'ls $WORKSPACE'
      }
    }
    stage('build') {
      steps {
        sh 'go build'
      }
    }
  }
}