import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import CIdField from '@/components/forms/fields/c-id-field.vue';

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'c-help-icon': true,
};

const snapshotStubs = {
  'c-help-icon': true,
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-id-field', () => {
  const factory = generateShallowRenderer(CIdField, { stubs });
  const snapshotFactory = generateRenderer(CIdField, {
    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

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

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
