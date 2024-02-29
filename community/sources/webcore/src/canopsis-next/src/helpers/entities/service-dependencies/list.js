import { uid } from '@/helpers/uid';

/**
 * @typedef {Service | Entity} ServiceDependency
 * @property {number} impact_state
 * @property {number} [depends_count]
 * @property {number} [state_depends_count]
 * @property {number} [impacts_count]
 */

/**
 * @typedef {ServiceDependency & ObjectKey} ServiceTreeviewDependency
 * @property {string} _id
 * @property {string[]} [children]
 */

/**
 * @typedef {Object & ObjectKey} LoadMoreTreeviewDependencyChild
 * @property {string} _id
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
 * Convert alarm to service dependency
 *
 * @param {Object} alarm
 * @returns {ServiceDependency}
 */
export const alarmToServiceDependency = alarm => ({
  ...alarm.entity,

  state: alarm?.v?.state?.val ?? 0,
  impact_state: alarm?.impact_state,
});

/**
 * Convert dependency to treeview dependency
 *
 * @param {ServiceDependency} [dependency = {}]
 * @param {boolean} [impact = false]
 * @returns {ServiceTreeviewDependency}
 */
export const dependencyToTreeviewDependency = (dependency = {}, impact = false) => {
  const preparedDependency = {
    entity: dependency,
    cycle: false,
    key: dependency._id,
    _id: dependency._id,
  };

  if ((!impact && dependency.depends_count) || (impact && dependency.impacts_count)) {
    preparedDependency.children = [];
  }

  return preparedDependency;
};

/**
 * Normalize treeview dependencies array
 *
 * @param {ServiceTreeviewDependency[]} [dependencies = []]
 * @param {boolean} [impact = false]
 * @returns {{result: string[], dependencies: Object}}
 */
export const normalizeDependencies = (dependencies = [], impact = false) => dependencies.reduce((acc, dependency) => {
  const preparedDependency = dependencyToTreeviewDependency(dependency, impact);

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
 * @param {string[]} [parentIds]
 * @returns {DenormalizedServiceTreeviewDependency[]}
 */
export const dependenciesDenormalize = ({
  ids = [],
  dependenciesByIds = {},
  metaByIds = {},
  parent,
  parentIds = [],
}) => ids
  .map((id) => {
    const meta = metaByIds[id] || {};
    const dependency = dependenciesByIds[id] || {};
    const { children, ...denormalizedDependency } = dependency;
    const isCycle = parentIds.includes(id);

    if (parent) {
      denormalizedDependency.parentKey = parent.key;
      denormalizedDependency.parentId = parent._id;
      denormalizedDependency.key = `${parent.key}/${denormalizedDependency.key}`;
    }

    denormalizedDependency.cycle = isCycle;

    if (!children || isCycle) {
      return denormalizedDependency;
    }

    denormalizedDependency.children = dependenciesDenormalize({
      ids: children,
      dependenciesByIds,
      metaByIds,
      parent: denormalizedDependency,
      parentIds: [...parentIds, id],
    });

    if (meta.page < meta.page_count && denormalizedDependency.children.length) {
      denormalizedDependency.children.push(getLoadMoreDenormalizedChild(denormalizedDependency));
    }

    return denormalizedDependency;
  });
