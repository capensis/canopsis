import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { PAGINATION_LIMIT } from '@/config';
import { FORM_CONTEXT_KEY } from '@/constants';

import CEntityCategoryField from '@/components/forms/fields/entity/c-entity-category-field.vue';

const { createEntityCategory, fetchCategoriesListWithoutStore } = vi.hoisted(() => ({
  createEntityCategory: vi.fn(),
  fetchCategoriesListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
}));

vi.mock('@/hooks/store/modules/entity-category', () => ({
  useEntityCategory: () => ({
    createEntityCategory,
    fetchCategoriesListWithoutStore,
  }),
}));

const lazySearchFieldBaseStub = createSelectInputStub('c-lazy-search-field');
const lazySearchFieldStub = {
  ...lazySearchFieldBaseStub,
  name: 'c-lazy-search-field',
  props: {
    ...lazySearchFieldBaseStub.props,
    clearable: Boolean,
  },
  emits: [
    'input',
    'change',
    'update:modelValue',
    'update:value',
    'update:search',
    'fetch',
    'fetch:more',
  ],
  template: `
    <div>
      <select
        class="c-lazy-search-field"
        :value="modelValue ?? value"
        @change="update"
      >
        <option v-for="item in availableItems" :value="getValue(item)" :key="getValue(item)">
          {{ item.value }}
        </option>
      </select>
      <slot name="prepend-item" />
    </div>
  `,
};

const stubs = {
  'c-lazy-search-field': lazySearchFieldStub,
  'c-help-icon': true,
};

const selectLazySearchField = wrapper => wrapper.findComponent(lazySearchFieldStub);

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
    vi.useRealTimers();
    createEntityCategory.mockReset();
    fetchCategoriesListWithoutStore.mockClear();
    fetchCategoriesListWithoutStore.mockResolvedValue({
      data: [],
      meta: { page_count: 1 },
    });
  });

  it('fetches the first categories page', async () => {
    factory();

    await flushPromises();

    expect(fetchCategoriesListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        limit: PAGINATION_LIMIT,
        search: '',
      }),
    });
  });

  it('hydrates a selected category by id', async () => {
    const [category] = entityCategories;
    fetchCategoriesListWithoutStore.mockImplementation(({ params }) =>
      Promise.resolve(
        params.ids
          ? { data: [category], meta: { page_count: 1 } }
          : { data: [], meta: { page_count: 1 } },
      ),
    );

    const wrapper = factory({
      propsData: {
        category: category._id,
      },
    });

    await flushPromises();

    expect(fetchCategoriesListWithoutStore).toHaveBeenCalledWith({
      params: {
        ids: [category._id],
        limit: 1,
      },
    });
    expect(selectLazySearchField(wrapper).props('value')).toEqual(category);
    expect(selectLazySearchField(wrapper).props('items')).toContainEqual(category);
  });

  it('fetches the next categories page', async () => {
    fetchCategoriesListWithoutStore.mockResolvedValue({
      data: entityCategories,
      meta: { page_count: 2 },
    });

    const wrapper = factory();

    await flushPromises();
    selectLazySearchField(wrapper).vm.$emit('fetch:more');
    await flushPromises();

    expect(fetchCategoriesListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 2,
        limit: PAGINATION_LIMIT,
      }),
    });
  });

  it('searches categories from the first page', async () => {
    vi.useFakeTimers();

    const wrapper = factory();

    await flushPromises();
    fetchCategoriesListWithoutStore.mockClear();

    selectLazySearchField(wrapper).vm.$emit('update:search', 'searched category');
    await vi.advanceTimersByTimeAsync(100);
    await flushPromises();

    expect(fetchCategoriesListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        search: 'searched category',
      }),
    });
  });

  it('emits the selected category object', async () => {
    fetchCategoriesListWithoutStore.mockResolvedValueOnce({
      data: entityCategories,
      meta: { page_count: 1 },
    });

    const wrapper = factory();

    await flushPromises();

    const [category] = entityCategories;

    selectLazySearchField(wrapper).vm.$emit('update:value', category);

    expect(wrapper).toEmitInput(category);
  });

  it('allows explicitly clearing a required category', () => {
    const wrapper = factory({
      propsData: {
        clearable: true,
        required: true,
      },
    });

    expect(selectLazySearchField(wrapper).props('clearable')).toBe(true);
  });

  it('renders the create category field only when categories are addable', () => {
    const wrapper = factory({
      propsData: {
        addable: true,
      },
    });

    const createField = wrapper.find('input');

    expect(createField.exists()).toBe(true);
    expect(createField.attributes('placeholder')).toBe('service.createCategory');
  });

  it('does not propagate create field mouse events to the category autocomplete', async () => {
    const wrapper = factory({
      attachTo: document.body,
      propsData: {
        addable: true,
      },
    });
    const autocompleteMousedownHandler = vi.fn();
    const createField = wrapper.find('input');

    wrapper.element.addEventListener('mousedown', autocompleteMousedownHandler);
    await createField.trigger('mousedown');
    createField.element.focus();

    expect(autocompleteMousedownHandler).not.toHaveBeenCalled();
    expect(document.activeElement).toBe(createField.element);

    wrapper.unmount();
  });

  it('creates and selects a category', async () => {
    const category = {
      _id: 'new-category-id',
      name: 'New category',
    };
    const form = {
      validateField: vi.fn().mockResolvedValue({ valid: true }),
      resetField: vi.fn(),
    };
    createEntityCategory.mockResolvedValue(category);

    const wrapper = factory({
      propsData: {
        addable: true,
      },
      parentComponent: {
        provide: {
          [FORM_CONTEXT_KEY]: form,
        },
      },
    });

    const createField = wrapper.find('input');

    await createField.setValue('  New category  ');
    await createField.trigger('keydown', { key: 'Enter' });
    await flushPromises();

    expect(createEntityCategory).toHaveBeenCalledWith({
      data: { name: category.name },
    });
    expect(wrapper).toEmitInput(category);
    expect(wrapper.vm.newCategory).toBe('');
  });
});
