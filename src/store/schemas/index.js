import { schema } from 'normalizr';

export const types = {
  ALARM: 'alarm',
  PBEHAVIOR: 'pbehavior',
  FILTER: 'filter',
  USER_PREFERENCES: 'userpreferences',
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

export const userPreferencesSchema = new schema.Entity(types.USER_PREFERENCES, {}, {
  idAttribute: '_id',
});

export default {
  [types.ALARM]: alarmSchema,
  [types.PBEHAVIOR]: pbehaviorSchema,
  [types.USER_PREFERENCES]: userPreferencesSchema,
};
