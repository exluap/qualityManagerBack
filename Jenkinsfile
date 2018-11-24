pipeline {
  agent {
    docker {
      image 'mhlg/rpi-golang'
    }

  }
  stages {
    stage('ready to build') {
      steps {
        sh 'go env'
      }
    }
    stage('build') {
      steps {
        sh 'ls $WORKSPACE'
        sh '''export GOPATH=/var/lib/jenkins/jobs/qualityManagerBack/branches/master/workspace

go get github.com/dgrijalva/jwt-go'''
      }
    }
  }
  environment {
    GIT_COMMITTER_NAME = 'exluapjenkins'
    GIT_COMMITTER_EMAIL = 'nickzaytsew+jenkins@gmail.com'
  }
}