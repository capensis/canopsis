import { createCheckboxInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import AvailabilityValueFilterMethodField from '@/components/other/availability/form/fields/availability-value-filter-method-field.vue';

const stubs = {
  'v-select': createCheckboxInputStub('v-select'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-select');

describe('availability-value-filter-method-field', () => {
  const factory = generateShallowRenderer(AvailabilityValueFilterMethodField, { stubs });

  test('Method changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_VALUE_FILTER_METHODS.greater,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('input', AVAILABILITY_VALUE_FILTER_METHODS.less);

    expect(wrapper).toEmitInput(AVAILABILITY_VALUE_FILTER_METHODS.less);
  });
});
