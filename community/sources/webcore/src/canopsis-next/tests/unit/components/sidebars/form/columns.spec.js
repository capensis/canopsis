import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

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
  const snapshotFactory = generateRenderer(Columns, { stubs });

  it('Columns changed after trigger columns field', () => {
    const wrapper = factory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        columns: [],
        label: '',
      },
    });

    const columnsField = selectColumnsField(wrapper);

    const columns = [{
      ...widgetColumnToForm(),

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
        type: ENTITIES_TYPES.alarm,
        label: 'Custom label',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `columns` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        type: ENTITIES_TYPES.alarm,
        label: 'Custom label',
        columns: [],
        withHtml: true,
        withColorIndicator: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
