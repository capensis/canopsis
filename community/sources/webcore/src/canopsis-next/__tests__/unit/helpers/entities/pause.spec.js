import { ENTITY_TYPES } from '@/constants';

import {
  entityPauseFormToRequest,
  entityResumeFormToRequest,
  isEntityPaused,
  isEntityPauseInherited,
  isPausePropagationAvailableForEntity,
  isResumeAvailableForEntity,
} from '@/helpers/entities/entity/pause';

describe('entity pause helpers', () => {
  test('detects own and inherited pauses', () => {
    const ownPauseEntity = { pause: { timestamp: 123 } };
    const inheritedPauseEntity = { pause: { timestamp: 123, inherited: true } };

    expect(isEntityPaused()).toBe(false);
    expect(isEntityPaused(ownPauseEntity)).toBe(true);
    expect(isEntityPauseInherited(ownPauseEntity)).toBe(false);
    expect(isEntityPauseInherited(inheritedPauseEntity)).toBe(true);
    expect(isResumeAvailableForEntity(ownPauseEntity)).toBe(true);
    expect(isResumeAvailableForEntity(inheritedPauseEntity)).toBe(false);
  });

  test('converts pause and resume forms to API payloads', () => {
    expect(entityPauseFormToRequest()).toEqual({ comment: '', type: '', inherited: false });
    expect(
      entityPauseFormToRequest({ comment: 'Comment', type: 'maintenance', inherited: true }),
    ).toEqual({ comment: 'Comment', type: 'maintenance', inherited: true });
    expect(entityResumeFormToRequest()).toEqual({ comment: '' });
    expect(entityResumeFormToRequest({ comment: 'Resume comment' })).toEqual({
      comment: 'Resume comment',
    });
  });

  test('allows pause propagation only for components and services', () => {
    expect(isPausePropagationAvailableForEntity({ type: ENTITY_TYPES.component })).toBe(true);
    expect(isPausePropagationAvailableForEntity({ type: ENTITY_TYPES.service })).toBe(true);
    expect(isPausePropagationAvailableForEntity({ type: ENTITY_TYPES.connector })).toBe(false);
    expect(isPausePropagationAvailableForEntity({ type: ENTITY_TYPES.resource })).toBe(false);
  });
});
