pipeline {
  agent any
  stages {
    stage('phase 1') {
      steps {
        parallel(
          "phase 1": {
            sh 'echo \'Hello World!!!\''
            echo 'Welcome to comtravo'
            
          },
          "phase 1.1": {
            echo 'phase 1.1'
            
          },
          "phase 1.2": {
            sh 'echo "foo bar"'
            
          }
        )
      }
    }
    stage('Phase 2') {
      steps {
        sh 'echo "in phase 2"'
      }
    }
  }
}