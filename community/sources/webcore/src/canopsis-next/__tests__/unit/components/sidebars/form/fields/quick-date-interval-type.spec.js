import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { QUICK_RANGES } from '@/constants';

import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/quick-date-interval-type.vue';

const stubs = {
  'widget-settings-item': true,
  'c-quick-date-interval-type-field': true,
};

const selectQuickDateIntervalTypeFieldField = wrapper =>
  wrapper.find('c-quick-date-interval-type-field-stub');

describe('quick-date-interval-type', () => {
  const factory = generateShallowRenderer(FieldQuickDateIntervalType, { stubs });

  test('Value changed after trigger number field', () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(QUICK_RANGES)).value;

    selectQuickDateIntervalTypeFieldField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
