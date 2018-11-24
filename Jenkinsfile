pipeline {
  agent {
    docker {
      image 'resin/raspberrypi3-golang'
      args '-v /tmp:/tmp'
    }

  }
  stages {
    stage('set worspace') {
      steps {
        sh 'export GOPATH=$WORKSPACE'
      }
    }
    stage('build') {
      steps {
        sh 'cat /etc/passwd'
        sh 'go get ./...'
      }
    }
  }
}