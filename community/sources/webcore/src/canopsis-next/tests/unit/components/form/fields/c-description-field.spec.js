import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CDescriptionField from '@/components/forms/fields/c-description-field.vue';

const stubs = {
  'v-textarea': true,
};

const selectTextareaNode = wrapper => wrapper.vm.$children[0];

describe('c-description-field', () => {
  const factory = generateShallowRenderer(CDescriptionField, { stubs });
  const snapshotFactory = generateRenderer(CDescriptionField, {
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Value changed after trigger the textarea', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = Faker.datatype.string();

    selectTextareaNode(wrapper).$emit('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Renders `c-description-field` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Default value',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-description-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        label: 'Custom label',
        name: 'customName',
        disabled: true,
        required: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-description-field` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '',
        required: true,
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});
