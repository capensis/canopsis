import { generateShallowRenderer } from '@unit/utils/vue';

import { SORT_ORDERS } from '@/constants';

import DefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';

const stubs = {
  'widget-settings-item': true,
};

const selectColumnSelectField = wrapper => wrapper.findAll('.c-select').at(0);
const selectOrderSelectField = wrapper => wrapper.findAll('.c-select').at(1);

describe('default-sort-column', () => {
  const columns = [
    { value: 'first', text: 'First' },
    { value: 'second', text: 'Second' },
  ];

  const factory = generateShallowRenderer(DefaultSortColumn, { stubs });

  it('Column changed after trigger select field with columns', () => {
    const [firstColumn, secondColumn] = columns;
    const wrapper = factory({
      propsData: {
        value: {
          column: firstColumn.value,
          order: SORT_ORDERS.desc,
        },
        columns,
      },
    });

    selectColumnSelectField(wrapper).triggerCustomEvent('input', secondColumn.value);

    expect(wrapper).toEmitInput({ column: secondColumn.value, order: SORT_ORDERS.desc });
  });

  it('Order changed after trigger select field with orders', () => {
    const [firstColumn] = columns;
    const wrapper = factory({
      propsData: {
        value: {
          column: firstColumn.value,
          order: SORT_ORDERS.desc,
        },
        columns,
      },
    });

    selectOrderSelectField(wrapper).triggerCustomEvent('input', SORT_ORDERS.asc);

    expect(wrapper).toEmitInput({ column: firstColumn.value, order: SORT_ORDERS.asc });
  });
});
