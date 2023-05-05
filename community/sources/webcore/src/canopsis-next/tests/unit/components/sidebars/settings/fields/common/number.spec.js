import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import FieldNumber from '@/components/sidebars/settings/fields/common/number.vue';

const stubs = {
  'widget-settings-item': true,
  'c-number-field': true,
};

const snapshotStubs = {
  'widget-settings-item': true,
  'c-number-field': true,
};

const selectNumberField = wrapper => wrapper.find('c-number-field-stub');

describe('field-number', () => {
  const factory = generateShallowRenderer(FieldNumber, {

    stubs,
  });

  const snapshotFactory = generateRenderer(FieldNumber, {

    stubs: snapshotStubs,
  });

  test('Value changed after trigger number field', () => {
    const wrapper = factory();

    const newValue = Faker.datatype.number();

    selectNumberField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `field-number` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `field-number` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 22,
        title: 'Custom title',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
