import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createInputStub } from '@unit/stubs/input';

import FastActionOutput from '@/components/sidebars/alarm/form/fields/fast-action-output.vue';

const stubs = {
  'widget-settings-item': true,
  'c-enabled-field': true,
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'widget-settings-item': true,
  'c-enabled-field': true,
};

const parentComponent = {
  provide: {
    list: {
      register: jest.fn(),
      unregister: jest.fn(),
    },
    listClick: jest.fn(),
  },
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('fast-action-output', () => {
  const factory = generateShallowRenderer(FastActionOutput, { stubs, parentComponent });
  const snapshotFactory = generateRenderer(FastActionOutput, { stubs: snapshotStubs, parentComponent });

  it('Enabled changed after trigger the enabled field', () => {
    const value = {
      enabled: true,
      value: '',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const enabledField = selectEnabledField(wrapper);

    enabledField.vm.$emit('input', false);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...value,
      enabled: false,
    });
  });

  it('Value changed after trigger the text field', () => {
    const value = {
      enabled: true,
      value: '',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.string();

    textField.setValue(newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...value,
      value: newValue,
    });
  });

  it('Renders `fast-action-output` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `fast-action-output` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          enabled: true,
          value: 'Value',
        },
        label: 'Label',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
