import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';

import EnabledLimit from '@/components/sidebars/settings/fields/common/enabled-limit.vue';

const localVue = createVueInstance();

const stubs = {
  'enabled-limit-field': true,
};

const snapshotFactory = (options = {}) => mount(EnabledLimit, {
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

const selectEnabledLimitField = wrapper => wrapper.find('enabled-limit-field-stub');

describe('enabled-limit', () => {
  it('Value changed after trigger enabled limit field', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `enabled-limit` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `enabled-limit` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          enabled: true,
          limit: 2,
        },
        title: 'Custom title',
        label: 'Custom label',
        optional: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `enabled-limit` with errors', async () => {
    const name = 'custom-name';
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          enabled: true,
          limit: 2,
        },
        title: 'Custom title',
        label: 'Custom label',
        optional: true,
      },
    });

    const validator = wrapper.getValidator();

    const columnsField = selectEnabledLimitField(wrapper);

    validator.attach({
      name,
      rules: 'required:true',
      getter: () => true,
      context: () => columnsField.vm,
      vm: columnsField.vm,
    });

    validator.errors.add({
      field: name,
      msg: 'error-message',
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
