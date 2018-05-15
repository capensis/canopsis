import { schema } from 'normalizr';

export const eventSchema = new schema.Entity('event', {}, { idAttribute: '_id' });

export default {
  event: eventSchema,
};
