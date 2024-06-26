package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ServiceMonitor(cr *v1alpha1.Keycloak) *monitoringv1.ServiceMonitor {
	return &monitoringv1.ServiceMonitor{
		ObjectMeta: v12.ObjectMeta{
			Name:      ServiceMonitorName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"monitoring-key": MonitoringKey,
			},
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				{
					Path:   "/auth/realms/master/metrics",
					Port:   "http",
					Scheme: "https",
				},
				{
					Path:   "/metrics",
					Port:   KeycloakMonitoringServiceName,
					Scheme: "http",
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": ApplicationName,
				},
			},
		},
	}
}

func ServiceMonitorSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ServiceMonitorName,
		Namespace: cr.Namespace,
	}
}
