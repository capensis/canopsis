import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { QUICK_RANGES } from '@/constants';

import CQuickDateIntervalTypeField from '@/components/forms/fields/compat/c-quick-date-interval-type-field.vue';

const stubs = {
  'c-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('c-quick-date-interval-type-field', () => {
  const factory = generateShallowRenderer(CQuickDateIntervalTypeField, { stubs });

  test('Value changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: QUICK_RANGES.last2Days,
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', QUICK_RANGES.yesterday);

    expect(wrapper).toEmitInput(QUICK_RANGES.yesterday);
  });
});
