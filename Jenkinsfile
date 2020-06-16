pipeline {
  agent any

  stages {
    stage("test")  {
      steps {
        docker version
        sh "echo v = $GOOS"
      }
    }
  }
}