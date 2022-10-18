import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import Columns from '@/components/sidebars/settings/fields/common/columns.vue';

const localVue = createVueInstance();

const stubs = {
  'c-columns-field': true,
};

const factory = (options = {}) => shallowMount(Columns, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Columns, {
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

const selectColumnsField = wrapper => wrapper.find('c-columns-field-stub');

describe('columns', () => {
  it('Columns changed after trigger columns field', () => {
    const wrapper = factory({
      propsData: {
        columns: [],
        label: '',
      },
    });

    const columnsField = selectColumnsField(wrapper);

    const columns = [{
      label: 'Column label',
      value: 'column.property',
      isHtml: false,
    }];

    columnsField.vm.$emit('input', columns);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(columns);
  });

  it('Renders `columns` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        label: 'Custom label',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `columns` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        label: 'Custom label',
        columns: [],
        withHtml: true,
        withColorIndicator: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
