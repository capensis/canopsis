import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';
import { SAMPLINGS } from '@/constants';

import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/sampling.vue';

const stubs = {
  'widget-settings-item': true,
  'c-sampling-field': true,
};

const snapshotStubs = {
  'widget-settings-item': true,
  'c-sampling-field': true,
};

const selectSamplingFieldField = wrapper => wrapper.find('c-sampling-field-stub');

describe('sampling', () => {
  const factory = generateShallowRenderer(FieldQuickDateIntervalType, { stubs });
  const snapshotFactory = generateRenderer(FieldQuickDateIntervalType, { stubs: snapshotStubs });

  test('Value changed after trigger sampling field', () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(SAMPLINGS));

    selectSamplingFieldField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `sampling` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `sampling` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SAMPLINGS.week,
        name: 'custom_name',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
