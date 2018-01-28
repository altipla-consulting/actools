
import altipla.CI


node {
  def ci = new CI()
  ci.project = 'altipla-tools'
  ci.zone = 'europe-west1-b'

  ci.init this

  stage('Checkout') {
    checkout scm
  }

  stage('dev-go') {
    ci.container name:'dev-go', context:'containers/dev-go'
  }

  stage('gcloud') {
    ci.container name:'gcloud', context:'containers/gcloud'
  }

  stage('glide') {
    ci.container name:'glide', context:'containers/glide'
  }

  stage('go') {
    ci.container name:'go', context:'containers/go'
  }

  stage('gulp') {
    ci.container name:'gulp', context:'containers/gulp'
  }

  stage('king') {
    ci.container name:'king', context:'containers/king'
  }

  stage('mysql') {
    ci.container name:'mysql', context:'containers/mysql'
  }

  stage('node') {
    ci.container name:'node', context:'containers/node'
  }

  stage('phpmyadmin') {
    ci.container name:'phpmyadmin', context:'containers/phpmyadmin'
  }

  stage('protoc') {
    ci.container name:'protoc', context:'containers/protoc'
  }

  stage('rambler') {
    ci.container name:'rambler', context:'containers/rambler'
  }

  stage('sass-lint') {
    ci.container name:'sass-lint', context:'containers/sass-lint'
  }

  stage('juice') {
    ci.container name:'juice', context:'containers/juice'
  }

  stage('actools') {
    sh 'actools go install ./cmd/actools'
    ci.gsutil "-h 'Cache-Control: no-cache' cp /root/.actools/cache-altipla-actools/bin/actools gs://altipla-tools-artifacts"
  }
}