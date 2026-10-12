import { ref } from 'vue';

import { hasVisibleFieldError } from '@/hooks/validator/validation-element-children';

describe('hasVisibleFieldError', () => {
  it('does not let an empty external errors prop mask vee-validate errors', () => {
    expect(hasVisibleFieldError({
      errors: ['The field is required'],
      instance: { errorMessages: [] },
    })).toBe(true);
  });

  it('detects errors exposed by a field component', () => {
    expect(hasVisibleFieldError({
      errors: [],
      instance: { $: { exposed: { errorMessages: ref(['Backend error']) } } },
    })).toBe(true);
  });

  it('returns false when neither error source contains a message', () => {
    expect(hasVisibleFieldError({
      errors: [],
      instance: { errorMessages: [] },
    })).toBe(false);
  });
});
