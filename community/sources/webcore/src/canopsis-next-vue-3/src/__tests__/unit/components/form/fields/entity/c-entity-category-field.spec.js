import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { MAX_LIMIT } from '@/constants';

import CEntityCategoryField from '@/components/forms/fields/entity/c-entity-category-field.vue';

const { fetchCategoriesListWithoutStore } = vi.hoisted(() => ({
  fetchCategoriesListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
  }),
}));

vi.mock('@/stores/modules/entity-category', () => ({
  useEntityCategoryStore: () => ({
    fetchCategoriesListWithoutStore,
  }),
}));

const stubs = {
  'c-select': createSelectInputStub('c-select'),
  'c-help-icon': true,
};

const selectSelectField = wrapper => wrapper.find('select.c-select');

const entityCategories = [
  {
    _id: 'c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf',
    name: 'Category',
    author: 'root',
    created: 1614861888,
    updated: 1614861888,
  },
  {
    _id: '441a2a17-9036-48a3-9ff7-f393487395a9',
    name: 'Category 2',
    author: 'root',
    created: 1614863990,
    updated: 1614863990,
  },
];

describe('c-entity-category-field', () => {
  const factory = generateShallowRenderer(CEntityCategoryField, { stubs });

  afterEach(() => {
    fetchCategoriesListWithoutStore.mockClear();
  });

  it('Check fetch list call', async () => {
    factory();

    await flushPromises();

    expect(fetchCategoriesListWithoutStore).toHaveBeenCalledWith({
      params: { limit: MAX_LIMIT },
    });
  });

  it('Category changed after trigger select field', async () => {
    fetchCategoriesListWithoutStore.mockResolvedValueOnce({
      data: entityCategories,
    });

    const wrapper = factory();

    await flushPromises();

    const [category] = entityCategories;

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', category);

    expect(wrapper).toEmitInput(category);
  });
});
