import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { SORT_ORDERS } from '@/constants';

import FieldDefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';

const stubs = {
  'widget-settings-item': true,
};

const selectColumnSelectField = wrapper => wrapper.findAll('.c-select').at(0);
const selectOrderSelectField = wrapper => wrapper.findAll('.c-select').at(1);

describe('field-sort-column', () => {
  const columns = [
    { value: 'first', text: 'First' },
    { value: 'second', text: 'Second' },
  ];

  const factory = generateShallowRenderer(FieldDefaultSortColumn, { stubs });

  test('Column changed after trigger combobox field', () => {
    const wrapper = factory({
      propsData: {
        columns,
      },
    });

    const newColumn = Faker.datatype.string();

    selectColumnSelectField(wrapper).triggerCustomEvent('input', newColumn);

    expect(wrapper).toEmitInput({
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
        columns,
      },
    });

    selectOrderSelectField(wrapper).triggerCustomEvent('input', SORT_ORDERS.desc);

    expect(wrapper).toEmitInput({
      column,
      order: SORT_ORDERS.desc,
    });
  });
});
