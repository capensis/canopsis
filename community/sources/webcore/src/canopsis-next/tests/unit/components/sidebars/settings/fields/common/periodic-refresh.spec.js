import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';

import PeriodicRefresh from '@/components/sidebars/settings/fields/common/periodic-refresh.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings-item': true,
  'periodic-refresh-field': true,
};

const factory = (options = {}) => shallowMount(PeriodicRefresh, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(PeriodicRefresh, {
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
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectPeriodicRefreshField = wrapper => wrapper.find('periodic-refresh-field-stub');

describe('periodic-refresh', () => {
  it('Unit as seconds settled created, if unit doesn\'t exist', () => {
    const wrapper = factory({
      propsData: {
        value: {
          value: 1,
        },
      },
    });

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      value: 1,
      unit: TIME_UNITS.second,
    });
  });

  it('Value changed after trigger periodic refresh field', () => {
    const wrapper = factory({
      propsData: {
        value: {
          value: 1,
          unit: TIME_UNITS.day,
        },
      },
    });

    const periodicRefreshField = selectPeriodicRefreshField(wrapper);

    const newValue = {
      value: 2,
      unit: TIME_UNITS.week,
    };

    periodicRefreshField.vm.$emit('input', newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  it('Renders `periodic-refresh` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `periodic-refresh` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          value: 1,
          unit: TIME_UNITS.minute,
        },
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `periodic-refresh` with errors', async () => {
    const name = 'custom-name';

    const wrapper = snapshotFactory({
      propsData: {
        value: {
          value: 1,
          unit: TIME_UNITS.minute,
        },
        name,
      },
    });

    const validator = wrapper.getValidator();

    const periodicRefreshField = selectPeriodicRefreshField(wrapper);

    validator.attach({
      name,
      rules: 'required:true',
      getter: () => true,
      context: () => periodicRefreshField.vm,
      vm: periodicRefreshField.vm,
    });

    validator.errors.add({
      field: name,
      msg: 'error-message',
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
