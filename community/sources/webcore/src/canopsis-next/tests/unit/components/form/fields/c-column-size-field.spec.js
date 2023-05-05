import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CColumnSizeField from '@/components/forms/fields/c-column-size-field.vue';

const stubs = {};

const factory = generateShallowRenderer(CColumnSizeField, { stubs,
});

const snapshotFactory = generateRenderer(CColumnSizeField, {

});

const selectRadioGroupNode = wrapper => wrapper.vm.$children[0];

describe('c-column-size-field', () => {
  test('Renders `c-column-size-field` with required props', () => {
    const wrapper = factory();

    const radioGroupNode = selectRadioGroupNode(wrapper);

    const newSize = Faker.datatype.number({ min: 1, max: 12 });

    radioGroupNode.$emit('change', newSize);

    expect(wrapper).toEmit('input', newSize);
  });

  test('Renders `c-column-size-field` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-column-size-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 12,
        name: 'custom_name',
        mobile: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-column-size-field` with tablet prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 4,
        tablet: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
