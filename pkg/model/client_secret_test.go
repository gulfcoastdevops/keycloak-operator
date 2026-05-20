package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientSecretCopiesTrueXPAnnotations(t *testing.T) {
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zta-sa",
			Namespace: "keycloak",
			Annotations: map[string]string{
				"truexp.truiem.io/project-service-account-credentials": "true",
				"truexp.truiem.io/tenant":                              "delta",
				"truexp.truiem.io/target-secret-name":                  "truexp-service-account-credentials",
				"example.com/ignore":                                   "ignored",
			},
		},
		Spec: v1alpha1.KeycloakClientSpec{
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "truexp-service-account",
				Secret:   "generated-secret",
			},
		},
	}

	secret := ClientSecret(cr)

	assert.Equal(t, "true", secret.Annotations["truexp.truiem.io/project-service-account-credentials"])
	assert.Equal(t, "delta", secret.Annotations["truexp.truiem.io/tenant"])
	assert.Equal(t, "truexp-service-account-credentials", secret.Annotations["truexp.truiem.io/target-secret-name"])
	assert.NotContains(t, secret.Annotations, "example.com/ignore")
}

func TestClientSecretReconciledPreservesNonTrueXPAnnotations(t *testing.T) {
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zta-sa",
			Namespace: "keycloak",
			Annotations: map[string]string{
				"truexp.truiem.io/tenant": "delta",
			},
		},
		Spec: v1alpha1.KeycloakClientSpec{
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "truexp-service-account",
				Secret:   "generated-secret",
			},
		},
	}
	current := &v1.Secret{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{
		"truexp.truiem.io/tenant": "old",
		"example.com/keep":        "kept",
	}}}

	reconciled := ClientSecretReconciled(cr, current)

	assert.Equal(t, "delta", reconciled.Annotations["truexp.truiem.io/tenant"])
	assert.Equal(t, "kept", reconciled.Annotations["example.com/keep"])
	assert.Equal(t, []byte("truexp-service-account"), reconciled.Data[ClientSecretClientIDProperty])
	assert.Equal(t, []byte("generated-secret"), reconciled.Data[ClientSecretClientSecretProperty])
}
