import { REQUEST_AUTH_TYPES, REQUEST_AUTH_TOKEN_TYPES } from '@/constants';

export default {
  parameterName: 'Parameter name',
  whereToAdd: 'Where to add',
  tokenValue: 'Token value',
  authTypes: {
    [REQUEST_AUTH_TYPES.none]: 'Not needed',
    [REQUEST_AUTH_TYPES.credentials]: 'With credentials',
    [REQUEST_AUTH_TYPES.token]: 'With token',
  },

  authTokenTypes: {
    [REQUEST_AUTH_TOKEN_TYPES.headerAuthorization]: 'Header (Authorization)',
    [REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer]: 'Header (Authorization Bearer)',
    [REQUEST_AUTH_TOKEN_TYPES.headerCustomParameter]: 'Header (Custom Parameter)',
    [REQUEST_AUTH_TOKEN_TYPES.payload]: '@:common.payload',
    [REQUEST_AUTH_TOKEN_TYPES.url]: 'URL',
  },
};
