import get from 'lodash/get';

/**
 * If entity has parent we should use this processStrategy
 *
 * @param entity
 * @param parent
 * @param key
 */
export const childProcessStrategy = (entity, parent, key) => {
  const result = { ...entity };

  if (parent && parent._embedded) {
    result._embedded = {
      parents: [{ type: parent._embedded.type, id: parent._id, key }],
    };
  }

  return result;
};

/**
 * If entity has parent we should use this mergeStrategy
 *
 * @param entityA
 * @param entityB
 */
export const childMergeStrategy = (entityA, entityB) => {
  const result = {
    ...entityA,
    ...entityB,
  };

  if (entityA._embedded || entityB._embedded) {
    result._embedded = {
      parents: [...get(entityA, '_embedded.parents', []), ...get(entityB, '_embedded.parents', [])],
    };
  }

  return result;
};

export default {
  childProcessStrategy,
  childMergeStrategy,
};

