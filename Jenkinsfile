pipeline {
    agent any

    environment {
        ARC_NAME = 'arcanhdeptrai'
        ARC_LOGIN_SERVER = 'arcanhdeptrai.azurecr.io'
        IMAGE_NAME = 'my-go-app'
        DOCKER_CREDENTIALS_ID = 'acr-demo'
    }
    stages {
        stage('ci:test') {
            steps {
                echo 'Running tests..'
                sh 'go test ./...'
            }
        }
        stage('ci:build') {
            steps {
                echo 'Building docker image..'
                script {
                    dockerImage.build("${ARC_LOGIN_SERVER}/${IMAGE_NAME}:latest")

                }
            }
        }
        stage('ci:push') {
            steps {
                echo 'Pushing docker image..'
                script {
                    docker.withRegistry("https://${ARC_LOGIN_SERVER}", DOCKER_CREDENTIALS_ID) {
                        dockerImage.push("latest")
                    }
                }
            }
        }
    }
    post {
        always {
            echo 'Cleaning up..'
            sh 'docker rmi ${ARC_LOGIN_SERVER}/${IMAGE_NAME}:latest || true'
        }
        failure {
            echo 'Build failed!'
        }
        success {
            echo 'Build succeeded and pushed to registry!'
        }
    }
}