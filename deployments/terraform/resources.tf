resource "kubernetes_namespace" "speedtest" {
  metadata {
    name = "speedtests"
  }
}

resource "helm_release" "influxdb" {
  name = "influxdb"
  repository = "https://helm.influxdata.com/"
  chart = "influxdb"
  values = [
    file("../helm/influxdb/my-values.yaml")
  ]
  namespace = kubernetes_namespace.speedtest.metadata[0].name
}

resource "helm_release" "speedtest" {
  name = "ookla-speedtest"
  chart = "../helm/ookla-speedtest"
  values = [
    file("../helm/ookla-speedtest/my-values.yaml")
  ]
  namespace = kubernetes_namespace.speedtest.metadata[0].name
}

resource "helm_release" "grafana" {
  name = "grafana"
  repository = "https://grafana.github.io/helm-charts"
  chart = "grafana" 
  values = [
    file("../helm/grafana/my-values.yaml")
  ]
  namespace = kubernetes_namespace.speedtest.metadata[0].name
}