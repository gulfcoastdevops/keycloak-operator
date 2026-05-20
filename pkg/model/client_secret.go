package model

import (
	"strings"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const trueXPAnnotationPrefix = "truexp.truiem.io/"

func clientSecretAnnotations(cr *v1alpha1.KeycloakClient) map[string]string {
	annotations := make(map[string]string)
	copyTrueXPAnnotations(annotations, cr.Annotations)
	if len(annotations) == 0 {
		return nil
	}
	return annotations
}

func clientSecretAnnotationsReconciled(cr *v1alpha1.KeycloakClient, currentAnnotations map[string]string) map[string]string {
	annotations := make(map[string]string)
	for key, value := range currentAnnotations {
		if !strings.HasPrefix(key, trueXPAnnotationPrefix) {
			annotations[key] = value
		}
	}
	copyTrueXPAnnotations(annotations, cr.Annotations)
	if len(annotations) == 0 {
		return nil
	}
	return annotations
}

func copyTrueXPAnnotations(dst, src map[string]string) {
	for key, value := range src {
		if strings.HasPrefix(key, trueXPAnnotationPrefix) {
			dst[key] = value
		}
	}
}

func ClientSecret(cr *v1alpha1.KeycloakClient) *v1.Secret {
	escapedSecretName := SanitizeResourceNameWithAlphaNum(ClientSecretName + "-" + cr.Name)
	return &v1.Secret{
		ObjectMeta: v12.ObjectMeta{
			Name:        escapedSecretName,
			Namespace:   cr.Namespace,
			Annotations: clientSecretAnnotations(cr),
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		Data: map[string][]byte{
			ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
			ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
		},
	}
}

func ClientSecretSelector(cr *v1alpha1.KeycloakClient) client.ObjectKey {
	escapedSelectorName := SanitizeResourceNameWithAlphaNum(ClientSecretName + "-" + cr.Name)
	return client.ObjectKey{
		Name:      escapedSelectorName,
		Namespace: cr.Namespace,
	}
}

func ClientSecretReconciled(cr *v1alpha1.KeycloakClient, currentState *v1.Secret) *v1.Secret {
	reconciled := currentState.DeepCopy()
	// Since the client is synced upon update, we always override what's there...
	reconciled.Annotations = clientSecretAnnotationsReconciled(cr, currentState.Annotations)
	reconciled.Data = map[string][]byte{
		ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
		ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
	}
	return reconciled
}

func DeprecatedClientSecret(cr *v1alpha1.KeycloakClient) *v1.Secret {
	escapedSecretName := SanitizeResourceNameWithAlphaNum(ClientSecretName + "-" + cr.Spec.Client.ClientID)
	return &v1.Secret{
		ObjectMeta: v12.ObjectMeta{
			Name:        escapedSecretName,
			Namespace:   cr.Namespace,
			Annotations: clientSecretAnnotations(cr),
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		Data: map[string][]byte{
			ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
			ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
		},
	}
}

func DeprecatedClientSecretSelector(cr *v1alpha1.KeycloakClient) client.ObjectKey {
	escapedSelectorName := SanitizeResourceNameWithAlphaNum(ClientSecretName + "-" + cr.Spec.Client.ClientID)
	return client.ObjectKey{
		Name:      escapedSelectorName,
		Namespace: cr.Namespace,
	}
}
