import { generateShallowRenderer } from '@unit/utils/vue';

import { SORT_ORDERS } from '@/constants';

import DefaultSortColumnsWithTemplate from '@/components/sidebars/form/fields/default-sort-columns-with-template.vue';

const selectSortColumnsForm = wrapper =>
  wrapper.findComponent({ name: 'widget-template-sort-columns-form' });

describe('default-sort-columns-with-template', () => {
  const templateId = 'template-id';
  const factory = generateShallowRenderer(DefaultSortColumnsWithTemplate);

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

  test('Applies an infos sort column using the matching current column value', () => {
    const wrapper = factory({
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

    wrapper.vm.updateTemplate(template);

    expect(wrapper.emitted('update:template')).toEqual([
      [
        templateId,
        [
          { sort_by: 'v.infos.rule.dictionary', sort: SORT_ORDERS.asc },
          { sort_by: 'v.component', sort: SORT_ORDERS.desc },
        ],
      ],
    ]);
  });

  test('Updates the displayed infos sort value when columns change', async () => {
    const sortColumn = { sort_by: 'v.infos', sort: SORT_ORDERS.asc };
    const wrapper = factory({
      propsData: {
        modelValue: [sortColumn],
        columns: [],
      },
    });

    expect(selectSortColumnsForm(wrapper).props('modelValue')).toEqual([sortColumn]);

    await wrapper.setProps({ columns: [{ value: 'v.infos.rule.dictionary' }] });

    expect(selectSortColumnsForm(wrapper).props('modelValue')).toEqual([
      { ...sortColumn, sort_by: 'v.infos.rule.dictionary' },
    ]);
  });
});
