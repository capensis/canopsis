import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createInputStub, createNumberInputStub } from '@unit/stubs/input';

import EnabledLimitField from '@/components/forms/fields/enabled-limit-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-enabled-field': createInputStub('c-enabled-field'),
  'c-number-field': createNumberInputStub('c-number-field'),
};

const snapshotStubs = {
  'c-enabled-field': true,
  'c-number-field': true,
};

const factory = (options = {}) => shallowMount(EnabledLimitField, {
  localVue,
  stubs,
  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(EnabledLimitField, {
  localVue,
  stubs: snapshotStubs,
  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectEnabledField = wrapper => wrapper.find('input.c-enabled-field');
const selectLimitField = wrapper => wrapper.find('input.c-number-field');

describe('enabled-limit-field', () => {
  it('Enabled changed after trigger the enabled field', () => {
    const value = {
      enabled: true,
      limit: 2,
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

  it('Limit changed after trigger the text field', () => {
    const value = {
      enabled: true,
      limit: 2,
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const limitField = selectLimitField(wrapper);

    const newLimit = 3;

    limitField.vm.$emit('input', newLimit);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...value,
      limit: newLimit,
    });
  });

  it('Renders `enabled-limit-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `enabled-limit-field` with custom props', () => {
    const value = {
      enabled: true,
      limit: 2,
    };
    const wrapper = snapshotFactory({
      propsData: {
        value,
        label: 'Custom label',
        name: 'custom-name',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
