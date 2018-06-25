/**
 * If entity has parent we should use this processStrategy
 *
 * @param entity
 * @param parent
 * @param key
 */
export const childProcessStrategy = (entity, parent, key) => ({
  ...entity,
  _embedded: {
    parents: [{ type: parent._embedded.type, id: parent._id, key }],
  },
});

/**
 * If entity has parent we should use this mergeStrategy
 *
 * @param entityA
 * @param entityB
 */
export const childMergeStrategy = (entityA, entityB) => ({
  ...entityA,
  ...entityB,
  _embedded: {
    parents: [...entityA._embedded.parents, ...entityB._embedded.parents],
  },
});

export default {
  childProcessStrategy,
  childMergeStrategy,
};
