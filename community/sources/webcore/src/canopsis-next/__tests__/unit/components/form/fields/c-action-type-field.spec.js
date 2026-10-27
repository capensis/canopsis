import { generateShallowRenderer } from '@unit/utils/vue';

import { ACTION_TYPES } from '@/constants';

import CActionTypeField from '@/components/forms/fields/c-action-type-field.vue';

describe('c-action-type-field', () => {
  const factory = generateShallowRenderer(CActionTypeField);

  it('uses Pause and Unpause action labels', () => {
    const wrapper = factory({
      propsData: {
        types: [ACTION_TYPES.pause, ACTION_TYPES.resume],
      },
    });

    expect(wrapper.vm.items).toEqual([
      { value: ACTION_TYPES.pause, text: 'Pause' },
      { value: ACTION_TYPES.resume, text: 'Unpause' },
    ]);
  });
});
