import { generateShallowRenderer } from '@unit/utils/vue';

import { ENTITIES_TYPES } from '@/constants';

import { widgetColumnToForm } from '@/helpers/entities/widget/column/form';

import Columns from '@/components/sidebars/form/fields/columns.vue';

const stubs = {
  'widget-settings-item': true,
  'c-columns-with-template-field': true,
};

const selectColumnsField = wrapper => wrapper.find('c-columns-with-template-field-stub');

describe('columns', () => {
  const factory = generateShallowRenderer(Columns, { stubs });

  it('Columns changed after trigger columns field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        columns: [],
        label: '',
      },
    });

    const columns = [
      {
        ...widgetColumnToForm(),

        label: 'Column label',
        value: 'column.property',
        isHtml: false,
      },
    ];

    selectColumnsField(wrapper).triggerCustomEvent('input', columns);

    expect(wrapper).toEmitInput(columns);
  });
});
