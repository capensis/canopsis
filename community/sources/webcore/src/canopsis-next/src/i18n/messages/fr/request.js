import { REQUEST_AUTH_TYPES, REQUEST_AUTH_TOKEN_TYPES } from '@/constants';

export default {
  parameterName: 'Nom du paramètre',
  whereToAdd: 'Où ajouter',
  tokenValue: 'Valeur du jeton',
  authTypes: {
    [REQUEST_AUTH_TYPES.none]: 'Non nécessaire',
    [REQUEST_AUTH_TYPES.credentials]: 'Avec identifiants',
    [REQUEST_AUTH_TYPES.token]: 'Avec jeton',
  },

  authTokenTypes: {
    [REQUEST_AUTH_TOKEN_TYPES.headerAuthorization]: 'En-tête (Authorization)',
    [REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer]: 'En-tête (Authorization Bearer)',
    [REQUEST_AUTH_TOKEN_TYPES.headerCustomParameter]: 'En-tête (paramètre personnalisé)',
    [REQUEST_AUTH_TOKEN_TYPES.payload]: '@:common.payload',
    [REQUEST_AUTH_TOKEN_TYPES.url]: 'URL',
  },
};
