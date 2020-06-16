pipeline {
  agent {
    dockerfile {
      filename "Dockerfile.dev"
    }
  }

  stages {
    stage("test")  {
      steps {
        docker version
        sh "echo v = $GOOS"
      }
    }
  }
}