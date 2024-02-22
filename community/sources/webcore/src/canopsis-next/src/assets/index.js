/**
 * @typedef {Function} RequireContextResult
 * @property {string | number} id
 * @property {Function} resolve
 * @property {Function} keys
 */

/**
 * Require all modules from context result
 *
 * @param {RequireContextResult} requireModule
 * @return {Object<string, unknown>}
 */
const requireAllModules = requireModule => requireModule.keys().reduce((acc, key) => {
  acc[key] = requireModule(key);

  return acc;
}, []);

export const eventFilterActionsTypesImages = requireAllModules(require.context('./event-filter-actions-types/'));
