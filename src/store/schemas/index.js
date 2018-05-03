import { schema } from 'normalizr';

export const eventSchema = new schema.Entity('event', {}, { idAttribute: 'id' });

export default {
  event: eventSchema,
};
