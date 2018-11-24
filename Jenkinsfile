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
        sh 'go get github.com/dgrijalva/jwt-go'
      }
    }
  }
}