import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CDescriptionField from '@/components/forms/fields/c-description-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-textarea': true,
};

const factory = (options = {}) => shallowMount(CDescriptionField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CDescriptionField, {
  localVue,

  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectTextareaNode = wrapper => wrapper.vm.$children[0];

describe('c-description-field', () => {
  test('Value changed after trigger the textarea', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const textareaNode = selectTextareaNode(wrapper);

    const newValue = Faker.datatype.string();

    textareaNode.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `c-description-field` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Default value',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
