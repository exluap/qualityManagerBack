pipeline {
  agent {
    docker {
      image '\'mhlg/rpi-golang\''
    }

  }
  stages {
    stage('ready to build') {
      steps {
        sh 'go get ./...'
      }
    }
    stage('build') {
      steps {
        sh 'go build'
      }
    }
  }
}