
import altipla.CI


node {
  def ci = new CI()
  ci.project = 'altipla-tools'
  ci.zone = 'europe-west1-b'

  ci.init this

  stage('Checkout') {
    checkout scm
  }

  def containers = []
  dir('containers') {
    findFiles().each { containers.push it.name }
  }

  containers.each {
    stage(it) {
      ci.container name:it, context:"containers/${it}"
    }
  }

  stage('actools') {
    sh 'actools go install ./cmd/actools'
    ci.gsutil "-h 'Cache-Control: no-cache' cp /root/.actools/cache-altipla-actools/bin/actools gs://tools.altipla.consulting/bin/actools"
  }

  stage('version manifest') {
    sh "echo ${ci.buildTag()} > version"
    ci.gsutil "-h 'Cache-Control: no-cache' cp version gs://tools.altipla.consulting/version-manifest/actools"
  }
}
