pipeline {
  agent {
    docker {
      image 'mhlg/rpi-golang'
    }

  }
  stages {
    stage('build') {
      steps {
        sh 'go build'
      }
    }
  }
}