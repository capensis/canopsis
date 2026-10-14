import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { SAMPLINGS } from '@/constants';

import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/sampling.vue';

const stubs = {
  'widget-settings-item': true,
  'c-sampling-field': true,
};

const selectSamplingFieldField = wrapper => wrapper.find('c-sampling-field-stub');

describe('sampling', () => {
  const factory = generateShallowRenderer(FieldQuickDateIntervalType, { stubs });

  test('Value changed after trigger sampling field', () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(SAMPLINGS));

    selectSamplingFieldField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
