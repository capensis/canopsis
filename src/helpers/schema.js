export const childProcessStrategy = (entity, parent, key) => ({
  ...entity,
  _embedded: {
    parents: [{ type: parent._embedded.type, id: parent._id, key }],
  },
});

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
