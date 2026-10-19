import { defineAsyncComponent } from 'vue';

/**
 * Convert legacy async component loaders from feature definitions into Vue 3 components.
 * Direct component objects are kept unchanged.
 *
 * @param {Object<string, Object|Function>} components
 * @returns {Object<string, Object>}
 */
export const normalizeFeatureComponents = (components = {}) =>
  Object.fromEntries(
    Object.entries(components).map(([name, component]) => [
      name,
      typeof component === 'function' ? defineAsyncComponent(component) : component,
    ]),
  );
