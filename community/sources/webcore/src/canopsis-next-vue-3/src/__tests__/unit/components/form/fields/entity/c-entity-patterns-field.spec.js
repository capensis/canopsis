import CEntityPatternsField from '@/components/forms/fields/entity/c-entity-patterns-field.vue';

describe('c-entity-patterns-field', () => {
  test('Multiple category values are displayed as closable chips', () => {
    const categories = [{ _id: 'category-id', name: 'Category' }];
    const categoryOptions = CEntityPatternsField.computed.categoryOptions.call({
      categories,
      categoriesPending: false,
    });

    const props = categoryOptions.valueField.props({ value: ['category-id'] });

    expect(props).toEqual(
      expect.objectContaining({
        items: categories,
        chips: true,
        multiple: true,
        closableChips: true,
      }),
    );
  });
});
