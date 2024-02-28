package common

import (
	"context"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RealmState struct {
	Realm            *kc.KeycloakRealm
	RealmUserSecrets map[string]*v1.Secret
	RealmClientScope map[string]*kc.KeycloakClientScope
	Context          context.Context
	Keycloak         *kc.Keycloak
}

func NewRealmState(context context.Context, keycloak kc.Keycloak) *RealmState {
	return &RealmState{
		Context:  context,
		Keycloak: &keycloak,
	}
}

func (i *RealmState) Read(cr *kc.KeycloakRealm, realmClient KeycloakInterface, controllerClient client.Client) error {
	realm, err := realmClient.GetRealm(cr.Spec.Realm.Realm)
	if err != nil {
		i.Realm = nil
		return err
	}

	i.Realm = realm
	if realm != nil && len(cr.Spec.Realm.Users) > 0 {
		// Get the state of the realm users
		i.RealmUserSecrets = make(map[string]*v1.Secret)
		for _, user := range cr.Spec.Realm.Users {
			secret, err := i.readRealmUserSecret(cr, user, controllerClient)
			if err != nil {
				return err
			}
			i.RealmUserSecrets[user.UserName] = secret

			cr.UpdateStatusSecondaryResources(SecretKind, model.GetRealmUserSecretName(i.Keycloak.Namespace, cr.Spec.Realm.Realm, user.UserName))
		}

	}

	// Get the state of realm scopes
	if realm != nil && len(cr.Spec.Realm.ClientScopes) > 0 {
		i.RealmClientScope = make(map[string]*kc.KeycloakClientScope)
		for _, scope := range cr.Spec.Realm.ClientScopes {
			s, err := i.readRealmClientScope(cr, scope, realmClient)
			if err != nil {
				return err
			}
			i.RealmClientScope[scope.Name] = s
		}
	}

	return nil
}

func (i *RealmState) readRealmUserSecret(realm *kc.KeycloakRealm, user *kc.KeycloakAPIUser, controllerClient client.Client) (*v1.Secret, error) {
	key := model.RealmCredentialSecretSelector(realm, user, i.Keycloak)
	secret := &v1.Secret{}

	// Try to find the user credential secret
	err := controllerClient.Get(i.Context, key, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return secret, err
}

func (i *RealmState) readRealmClientScope(cr *kc.KeycloakRealm, scope kc.KeycloakClientScope, realmClient KeycloakInterface) (*kc.KeycloakClientScope, error) {
	// Get the state of the realm client scopes
	clientScope, err := realmClient.ListAvailableClientScopes(cr.Spec.Realm.Realm)
	if err != nil {
		return nil, err
	}

	// Find the client scope
	for _, s := range clientScope {
		if s.Name == scope.Name {
			return &s, nil
		}
	}

	return nil, nil
}
