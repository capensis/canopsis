import { convertContextUserPreferenceToQuery } from '@/helpers/entities/entity/query';

describe('context entity query helpers', () => {
  it('disables idle entities when the preference is not set', () => {
    expect(convertContextUserPreferenceToQuery()).toEqual({
      category: undefined,
      filter: undefined,
      no_events: false,
    });
  });

  it('keeps the idle entities preference', () => {
    expect(convertContextUserPreferenceToQuery({
      content: { noEvents: true },
    }).no_events).toBe(true);
  });
});
