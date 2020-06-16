pipeline {
  agent {
    dockerfile {
      filename "Dockerfile.dev"
    }
  }

  environment {

  }

  stages {
    stage("test")  {
      steps {
        sh "echo v = $GOOS"
      }
    }
  }
}