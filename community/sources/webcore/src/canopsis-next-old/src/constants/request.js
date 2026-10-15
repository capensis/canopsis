export const REQUEST_AUTH_TYPES = {
  none: 'none',
  credentials: 'credentials',
  token: 'token',
};

export const REQUEST_AUTH_TOKEN_TYPES = {
  headerAuthorization: 'headerAuthorization',
  headerAuthorizationBearer: 'headerAuthorizationBearer',
  headerCustomParameter: 'headerCustomParameter',
  payload: 'payload',
  url: 'url',
};

export const REQUEST_AUTH_TOKEN_TYPES_TO_HEADERS = {
  [REQUEST_AUTH_TOKEN_TYPES.headerAuthorization]: 'Authorization',
  [REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer]: 'Authorization: Bearer',
};

export const HEADERS_TO_REQUEST_AUTH_TOKEN_TYPES = Object.fromEntries(
  Object.entries(REQUEST_AUTH_TOKEN_TYPES_TO_HEADERS).map(([key, value]) => [value, key]),
);
