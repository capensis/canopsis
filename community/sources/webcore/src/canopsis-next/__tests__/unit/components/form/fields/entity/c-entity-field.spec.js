import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { PAGINATION_LIMIT } from '@/config';
import { BASIC_ENTITY_TYPES } from '@/constants';

import CEntityField from '@/components/forms/fields/compat/c-entity-field.vue';

const { fetchContextEntitiesListWithoutStore } = vi.hoisted(() => ({
  fetchContextEntitiesListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
}));

vi.mock('@/hooks/store/modules/entity', () => ({
  useEntity: () => ({
    fetchContextEntitiesListWithoutStore,
  }),
}));

const stubs = {
  'c-lazy-search-field': createSelectInputStub('c-lazy-search-field'),
};

const selectLazySearchField = wrapper => wrapper.find('select.c-lazy-search-field');

describe('c-entity-field', () => {
  const items = [
    {
      _id: 'id',
      value: 'value',
      text: 'Text',
      type: 'type1',
    },
    {
      _id: 'id 2',
      value: 'value 2',
      text: 'Text 2',
      type: 'type2',
    },
    {
      _id: 'id 3',
      value: 'value 3',
      text: 'Text 3',
      type: 'type3',
    },
  ];

  const factory = generateShallowRenderer(CEntityField, { stubs });

  afterEach(() => {
    fetchContextEntitiesListWithoutStore.mockClear();
  });

  it('Entities fetched after mount', async () => {
    factory({
      propsData: {
        itemText: 'text',
        itemValue: 'value',
      },
    });

    await flushPromises();

    expect(fetchContextEntitiesListWithoutStore).toHaveBeenCalledWith({
      params: {
        limit: PAGINATION_LIMIT,
        page: 1,
        search: '',
        type: Object.values(BASIC_ENTITY_TYPES),
      },
    });
  });

  it('Value changed after trigger the input', async () => {
    const wrapper = factory({
      propsData: {
        itemText: 'text',
        itemValue: 'value',
      },
    });

    await flushPromises();

    selectLazySearchField(wrapper).triggerCustomEvent('update:value', items[0].value);

    expect(wrapper).toEmitInput(items[0].value);
  });
});
