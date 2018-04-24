import { schema } from 'normalizr';

export const alarmSchema = new schema.Entity('alarm', {}, { idAttribute: '_id' });

export default {
  alarm: alarmSchema,
};
