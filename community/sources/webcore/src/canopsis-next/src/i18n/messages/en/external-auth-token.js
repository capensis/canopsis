export default {
  tokenName: 'Token name',
  tokenExpirationTime: 'Token expiration time',
  lastUsedDate: 'Last used date',
  lastUpdateDate: 'Last update date',
  allowVariables: 'Allow variables in the Token',
  tokenExpirationHelpText: 'API Response field in JSON format where “Token” is taken from',
  tokenCanNotBeDeleted: 'Token can not be deleted, it is used in:\n{rules}',
  urlHelp: '<p>The accessible variables are: <strong>.Env</strong></p>'
    + '<i>For example:</i>'
    + '<pre>"https://exampleurl.com?env={{ .Env.System.ENV_var }}"</pre>',
};
