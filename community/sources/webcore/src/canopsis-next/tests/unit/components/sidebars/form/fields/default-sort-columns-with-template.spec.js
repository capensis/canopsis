import { generateShallowRenderer } from '@unit/utils/vue';

import { CUSTOM_WIDGET_TEMPLATE, SORT_ORDERS } from '@/constants';

import DefaultSortColumnsWithTemplate from '@/components/sidebars/form/fields/default-sort-columns-with-template.vue';

const selectSortColumnsForm = wrapper => wrapper.findComponent({ name: 'widget-template-sort-columns-form' });
const selectTemplateField = wrapper => wrapper.find('c-widget-template-field-stub');

describe('default-sort-columns-with-template', () => {
  const templateId = 'template-id';
  const factory = generateShallowRenderer(DefaultSortColumnsWithTemplate, {
    stubs: ['c-widget-template-field', 'c-alert'],
  });

  test('Does not mark an infos sort column as missing when a column starts with it', () => {
    const wrapper = factory({
      propsData: {
        template: templateId,
        templates: [
          {
            _id: templateId,
            sort_columns: [
              { sort_by: 'v.infos', sort: SORT_ORDERS.asc },
              { sort_by: 'v.component', sort: SORT_ORDERS.desc },
            ],
          },
        ],
        columns: [
          { value: 'v.infos.rule.dictionary' },
          { value: 'v.component-suffix' },
        ],
      },
    });

    expect(wrapper.vm.missingColumns).toHaveLength(1);
    expect(wrapper.vm.missingColumns[0].sort_by).toBe('v.component');
  });

  test('Applies an infos sort column using the matching current column value', async () => {
    const updateTemplate = jest.fn();
    const wrapper = factory({
      listeners: { 'update:template': updateTemplate },
      propsData: {
        columns: [
          { value: 'v.infos.rule.dictionary' },
          { value: 'v.component' },
        ],
      },
    });
    const template = {
      value: templateId,
      sort_columns: [
        { sort_by: 'v.infos', sort: SORT_ORDERS.asc },
        { sort_by: 'v.component', sort: SORT_ORDERS.desc },
      ],
    };

    await selectTemplateField(wrapper).triggerCustomEvent('input', template);

    expect(updateTemplate).toHaveBeenCalledTimes(1);
    expect(updateTemplate).toHaveBeenCalledWith(
      templateId,
      [
        { sort_by: 'v.infos.rule.dictionary', sort: SORT_ORDERS.asc },
        { sort_by: 'v.component', sort: SORT_ORDERS.desc },
      ],
    );
    expect(template.sort_columns[0].sort_by).toBe('v.infos');
  });

  test('Updates the displayed infos sort value when columns change', async () => {
    const sortColumn = { sort_by: 'v.infos', sort: SORT_ORDERS.asc };
    const wrapper = factory({
      propsData: {
        sortColumns: [sortColumn],
        columns: [],
      },
    });

    expect(selectSortColumnsForm(wrapper).props('value')).toEqual([sortColumn]);

    await wrapper.setProps({ columns: [{ value: 'v.infos.rule.dictionary' }] });

    expect(selectSortColumnsForm(wrapper).props('value')).toEqual([
      { ...sortColumn, sort_by: 'v.infos.rule.dictionary' },
    ]);
    expect(sortColumn.sort_by).toBe('v.infos');
  });

  test('Clears the selected template', async () => {
    const updateTemplate = jest.fn();
    const wrapper = factory({
      propsData: { template: templateId },
      listeners: { 'update:template': updateTemplate },
    });

    await selectTemplateField(wrapper).triggerCustomEvent('input', null);

    expect(updateTemplate).toHaveBeenCalledTimes(1);
    expect(updateTemplate).toHaveBeenCalledWith(undefined, undefined);
  });

  test('Emits input when editing custom sort columns', async () => {
    const wrapper = factory({ propsData: { template: CUSTOM_WIDGET_TEMPLATE } });
    const sortColumns = [{ sort_by: 'v.component', sort: SORT_ORDERS.asc }];

    await selectSortColumnsForm(wrapper).triggerCustomEvent('input', sortColumns);

    expect(wrapper).toEmitInput(sortColumns);
  });
});
