pipeline {
  agent {
    docker {
      image 'library/golang'
    }

  }
  stages {
    stage('get depend') {
      steps {
        sh 'go get ./...'
      }
    }
    stage('build') {
      parallel {
        stage('build') {
          steps {
            sh 'go build'
          }
        }
        stage('error') {
          steps {
            archiveArtifacts(artifacts: '*', excludes: '*.go, *.db', onlyIfSuccessful: true, defaultExcludes: true)
          }
        }
      }
    }
  }
}