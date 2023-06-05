import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import CIdField from '@/components/forms/fields/c-id-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'c-help-icon': true,
};

const snapshotStubs = {
  'c-help-icon': true,
};

const factory = (options = {}) => shallowMount(CIdField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CIdField, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-id-field', () => {
  test('Value changed after trigger the text field', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.string();

    textField.setValue(newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `c-id-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Default value',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-id-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'Custom value',
        label: 'Custom label',
        name: 'customName',
        helpText: 'helpText',
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-id-field` with errors', async () => {
    const name = 'customName';

    const wrapper = snapshotFactory({
      propsData: {
        value: '',
        name,
      },
    });

    const validator = wrapper.getValidator();

    validator.errors.add([
      {
        field: name,
        msg: 'Value error',
      },
    ]);

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});
