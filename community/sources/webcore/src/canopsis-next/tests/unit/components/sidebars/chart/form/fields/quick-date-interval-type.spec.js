import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';
import { QUICK_RANGES } from '@/constants';

import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/quick-date-interval-type.vue';

const stubs = {
  'widget-settings-item': true,
  'c-quick-date-interval-type-field': true,
};

const snapshotStubs = {
  'widget-settings-item': true,
  'c-quick-date-interval-type-field': true,
};

const selectQuickDateIntervalTypeFieldField = wrapper => wrapper.find('c-quick-date-interval-type-field-stub');

describe('quick-date-interval-type', () => {
  const factory = generateShallowRenderer(FieldQuickDateIntervalType, { stubs });
  const snapshotFactory = generateRenderer(FieldQuickDateIntervalType, { stubs: snapshotStubs });

  test('Value changed after trigger number field', () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(QUICK_RANGES));

    selectQuickDateIntervalTypeFieldField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `quick-date-interval-type` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `quick-date-interval-type` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: QUICK_RANGES.last1Year.value,
        name: 'custom_name',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
