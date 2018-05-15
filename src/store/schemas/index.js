import { schema } from 'normalizr';

export const pbehaviorSchema = new schema.Entity('pbehavior', {}, {
  idAttribute: '_id',
  processStrategy: (entity, parent, key) =>
    ({
      ...entity,
      _embedded: {
        parentId: parent._id,
        parentType: parent._embedded.type,
        relationType: key,
      },
    }),
});

export const alarmSchema = new schema.Entity('alarm', {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
  processStrategy: entity => ({ ...entity, _embedded: { type: 'alarm' } }),
});

export default {
  alarm: alarmSchema,
  pbehavior: pbehaviorSchema,
};
