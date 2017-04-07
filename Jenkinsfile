pipeline {
  agent any
  stages {
    stage('libdeps') {
      steps {
        sh 'make libdeps'
      }
    }
    stage('all') {
      steps {
        sh 'make all'
      }
    }
  }
}