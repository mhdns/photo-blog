pipeline {
  agent {
    dockerfile {
      filename "Dockerfile.dev"
    }
  }

  stages {
    stage("test")  {
      steps {
        sh "echo v = $GOOS"
      }
    }
  }
}