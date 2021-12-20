import { get, omit, isUndefined } from 'lodash';

import uid from '../uid';

/**
 * @typedef {Object} ServiceDependency
 * @property {Service | Entity} entity
 * @property {Object} [alarm = null]
 * @property {number} impact_state
 * @property {boolean} has_dependencies
 * @property {boolean} has_impacts
 */

/**
 * @typedef {ServiceDependency} ServiceTreeviewDependency
 * @property {string} _id
 * @property {string} key
 * @property {string[]} [children]
 */

/**
 * @typedef {Object} LoadMoreTreeviewDependencyChild
 * @property {string} _id
 * @property {string} key
 * @property {string} [parentKey]
 * @property {string} [parentId]
 */

/**
 * @typedef {
 *   DenormalizedServiceTreeviewDependency
 *   | LoadMoreTreeviewDependencyChild
 * } DenormalizedServiceTreeviewDependencyChild
 */

/**
 * @typedef {ServiceTreeviewDependency} DenormalizedServiceTreeviewDependency
 * @property {DenormalizedServiceTreeviewDependency[]} [children]
 */

/**
 * Convert service to service dependency
 *
 * @param {Service} entity
 * @param {Object} [alarm = null]
 * @returns {ServiceDependency}
 */
export const serviceToServiceDependency = (entity, alarm = null) => {
  let impactState = get(entity, 'impact_state');

  if (isUndefined(impactState)) {
    impactState = get(alarm, 'impact_state', 0);
  }

  return {
    entity,
    alarm,

    impact_state: impactState,
    has_dependencies: true,
  };
};

/**
 * Convert dependency to treeview dependency
 *
 * @param {ServiceDependency} [dependency = {}]
 * @returns {ServiceTreeviewDependency}
 */
export const dependencyToTreeviewDependency = (dependency = {}) => {
  const preparedDependency = {
    ...dependency,

    entity: { ...dependency.entity, impact_state: dependency.impact_state },
    key: uid('dependency'),
    _id: dependency.entity._id,
  };

  if (dependency.has_dependencies || dependency.has_impacts) {
    preparedDependency.children = [];
  }

  return preparedDependency;
};

/**
 * Convert treeview dependency to dependency
 *
 * @param {ServiceTreeviewDependency} [treeviewDependency = {}]
 * @returns {ServiceDependency}
 */
export const treeviewDependencyToDependency = (treeviewDependency = {}) => omit(
  treeviewDependency,
  ['children', 'parentId', 'parentKey'],
);

/**
 * Normalize treeview dependencies array
 *
 * @param {ServiceTreeviewDependency[]} [dependencies = []]
 * @returns {{result: string[], dependencies: Object}}
 */
export const normalizeDependencies = (dependencies = []) => dependencies.reduce((acc, dependency) => {
  const preparedDependency = dependencyToTreeviewDependency(dependency);

  acc.dependencies[preparedDependency._id] = preparedDependency;
  acc.result.push(preparedDependency._id);

  return acc;
}, { result: [], dependencies: {} });

/**
 * Get child item which will show load more button
 *
 * @param {ServiceTreeviewDependency} parent
 * @returns {LoadMoreTreeviewDependencyChild}
 */
export const getLoadMoreDenormalizedChild = (parent) => {
  const key = uid('dependency');
  const child = {
    key,

    _id: key,
    loadMore: true,
  };

  if (parent) {
    child.parentId = parent._id;
    child.parentKey = parent.key;
  }

  return child;
};

/**
 * Denormalize treeview dependencies array
 *
 * @param {string[]} [ids = []]
 * @param {Object} [dependenciesByIds = {}]
 * @param {Object} [metaByIds = {}]
 * @param {DenormalizedServiceTreeviewDependency} [parent]
 * @returns {DenormalizedServiceTreeviewDependency[]}
 */
export const dependenciesDenormalize = (ids = [], dependenciesByIds = {}, metaByIds = {}, parent) => ids
  .map((id) => {
    const meta = metaByIds[id] || {};
    const dependency = dependenciesByIds[id] || {};
    const { children } = dependency;
    const denormalizedDependency = { ...dependency };

    if (parent) {
      denormalizedDependency.parentKey = parent.key;
      denormalizedDependency.parentId = parent._id;
    }

    if (!children) {
      return denormalizedDependency;
    }

    denormalizedDependency.children = dependenciesDenormalize(
      children,
      dependenciesByIds,
      metaByIds,
      denormalizedDependency,
    );

    if (meta.page < meta.page_count) {
      denormalizedDependency.children.push(getLoadMoreDenormalizedChild(dependency));
    }

    return denormalizedDependency;
  });
