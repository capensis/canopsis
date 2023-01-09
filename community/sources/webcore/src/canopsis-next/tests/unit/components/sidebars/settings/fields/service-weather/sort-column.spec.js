import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { SORT_ORDERS } from '@/constants';

import FieldSortColumn from '@/components/sidebars/settings/fields/service-weather/sort-column.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings-item': true,
  'v-select': createSelectInputStub('v-select'),
  'v-combobox': createSelectInputStub('v-combobox'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectColumnField = wrapper => wrapper.find('.v-combobox');
const selectOrderField = wrapper => wrapper.find('.v-select');

describe('field-sort-column', () => {
  const factory = generateShallowRenderer(FieldSortColumn, {
    localVue,
    stubs,
  });

  const snapshotFactory = generateRenderer(FieldSortColumn, {
    localVue,
    stubs: snapshotStubs,
  });

  test('Column changed after trigger combobox field', () => {
    const wrapper = factory();

    const newColumn = Faker.datatype.string();

    selectColumnField(wrapper).vm.$emit('input', newColumn);

    expect(wrapper).toEmit('input', {
      column: newColumn,
      order: SORT_ORDERS.asc,
    });
  });

  test('Order changed after trigger select field', () => {
    const column = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        value: {
          column,
          order: SORT_ORDERS.asc,
        },
      },
    });

    selectOrderField(wrapper).vm.$emit('input', SORT_ORDERS.desc);

    expect(wrapper).toEmit('input', {
      column,
      order: SORT_ORDERS.desc,
    });
  });

  test('Renders `field-sort-column` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `field-sort-column` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          column: 'sort_column',
          order: SORT_ORDERS.desc,
        },
        columns: [{
          value: 'column',
          label: 'Column',
        }],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
