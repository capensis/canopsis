import { schema } from 'normalizr';

export const types = {
  ALARM: 'alarm',
  PBEHAVIOR: 'pbehavior',
  USER_PREFERENCE: 'user_preference',
};

export const pbehaviorSchema = new schema.Entity(types.PBEHAVIOR, {}, {
  idAttribute: '_id',
  processStrategy: (entity, parent, key) =>
    ({
      ...entity,
      _embedded: {
        parentType: types.ALARM,
        parentId: parent._id,
        relationType: key,
      },
    }),
});

export const alarmSchema = new schema.Entity(types.ALARM, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
});

export const userPreferenceSchema = new schema.Entity(types.USER_PREFERENCE, {}, {
  idAttribute: '_id',
});

export default {
  [types.ALARM]: alarmSchema,
  [types.PBEHAVIOR]: pbehaviorSchema,
  [types.USER_PREFERENCE]: userPreferenceSchema,
};
