resource "kubernetes_namespace" "speedtest" {
  metadata {
    name = "speedtest"
  }
}

resource "helm_release" "speedtest" {
  name = "ookla-speedtest"
  chart = "../helm/ookla-speedtest"
  values = [
    file("../helm/ookla-speedtest/my-values.yaml")
  ]
  namespace = kubernetes_namespace.speedtest.metadata[0].name
}