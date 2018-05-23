import { schema } from 'normalizr';

export const pbehaviorSchema = new schema.Entity('pbehavior', {}, { idAttribute: '_id' });

export const alarmSchema = new schema.Entity('alarm', {
  pbehaviors: [pbehaviorSchema],
}, { idAttribute: '_id' });


export default {
  alarm: alarmSchema,
  pbehavior: pbehaviorSchema,

};
