import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { EVENT_FILTER_TYPES } from '@/constants';

import CEventFilterTypeField from '@/components/forms/fields/c-event-filter-type-field.vue';

const stubs = {
  'c-select': createSelectInputStub('c-select'),
};

const selectSelectField = wrapper => wrapper.find('select.c-select');

describe('c-event-filter-type-field', () => {
  const factory = generateShallowRenderer(CEventFilterTypeField, { stubs });

  test('Value changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        modelValue: '',
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', EVENT_FILTER_TYPES.drop);

    expect(wrapper).toEmitInput(EVENT_FILTER_TYPES.drop);
  });
});
