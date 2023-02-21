import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { ENTITIES_TYPES } from '@/constants';

import { widgetColumnToForm } from '@/helpers/forms/shared/widget-column';

import Columns from '@/components/sidebars/settings/fields/common/columns.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings-item': true,
  'c-columns-with-template-field': true,
};

const factory = (options = {}) => shallowMount(Columns, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Columns, {
  localVue,
  stubs,

  ...options,
});

const selectColumnsField = wrapper => wrapper.find('c-columns-with-template-field-stub');

describe('columns', () => {
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
