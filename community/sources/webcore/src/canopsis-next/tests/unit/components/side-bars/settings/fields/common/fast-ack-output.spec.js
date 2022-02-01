import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createInputStub } from '@unit/stubs/input';
import FastAckOutput from '@/components/sidebars/settings/fields/alarm/fast-ack-output.vue';

const localVue = createVueInstance();

const stubs = {
  'c-enabled-field': true,
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'c-enabled-field': true,
};

const factory = (options = {}) => shallowMount(FastAckOutput, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      list: {
        register: jest.fn(),
        unregister: jest.fn(),
      },
      listClick: jest.fn(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(FastAckOutput, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    provide: {
      list: {
        register: jest.fn(),
        unregister: jest.fn(),
      },
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('fast-ack-output', () => {
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

  it('Renders `fast-ack-output` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `fast-ack-output` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          enabled: true,
          value: 'Value',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
