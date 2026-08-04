pipeline {
  agent any

  environment {
    IMAGE = "3224142123/cv-backend"
    DEPLOY = "cv-backend"
    NS = "cv-portfolio"
  }

  stages {
    stage('Build') {
      steps {
        script {
          docker.build("${IMAGE}:${GIT_COMMIT}", "-f Dockerfile.prod .")
        }
      }
    }

    stage('Push') {
      steps {
        script {
          docker.withRegistry('https://index.docker.io/v1/', 'dockerhub-credentials') {
            docker.image("${IMAGE}:${GIT_COMMIT}").push()
            docker.image("${IMAGE}:${GIT_COMMIT}").push('latest')
          }
        }
      }
    }

    stage('Deploy') {
      steps {
        withKubeConfig([credentialsId: 'kubeconfig']) {
          sh """
            kubectl set image deployment/${DEPLOY} \
              ${DEPLOY}=${IMAGE}:${GIT_COMMIT} \
              -n ${NS}
            kubectl rollout status deployment/${DEPLOY} -n ${NS}
          """
        }
      }
    }

    stage('Verify') {
      steps {
        withKubeConfig([credentialsId: 'kubeconfig']) {
          sh '''
            kubectl run cv-backend-verify --rm -i --restart=Never \
              --image=curlimages/curl -n ${NS} -- \
              curl -sf http://cv-backend:3333/healthz | grep -q '"status":"ok"'
          '''
        }
      }
    }
  }

  post {
    failure {
      echo 'Pipeline failed!'
    }
  }
}
