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
    stage('build') {
      steps {
        sh 'cat /etc/passwd'
        sh 'go get ./...'
      }
    }
  }
}