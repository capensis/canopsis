import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { PAGINATION_LIMIT } from '@/config';

import CPbehaviorTypeField from '@/components/forms/fields/pbehavior/c-pbehavior-type-field.vue';

const { fetchPbehaviorTypesFieldList } = vi.hoisted(() => ({
  fetchPbehaviorTypesFieldList: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
}));

vi.mock('@/hooks/store/modules/pbehavior-type', () => ({
  usePbehaviorType: () => ({
    fetchPbehaviorTypesFieldList,
  }),
}));

const lazySearchFieldBaseStub = createSelectInputStub('c-lazy-search-field');
const lazySearchFieldStub = {
  ...lazySearchFieldBaseStub,
  name: 'c-lazy-search-field',
  emits: [...lazySearchFieldBaseStub.emits, 'fetch', 'fetch:more', 'update:search'],
};

const stubs = {
  'c-lazy-search-field': lazySearchFieldStub,
};

const pbehaviorTypes = [
  { _id: 'type-1', name: 'Type 1' },
  { _id: 'type-2', name: 'Type 2' },
];

describe('c-pbehavior-type-field', () => {
  const factory = generateShallowRenderer(CPbehaviorTypeField, { stubs });

  afterEach(() => {
    vi.useRealTimers();
    fetchPbehaviorTypesFieldList.mockClear();
    fetchPbehaviorTypesFieldList.mockResolvedValue({
      data: [],
      meta: { page_count: 1 },
    });
  });

  test('Fetches the first page of requested PBehavior types', async () => {
    const types = ['pause'];

    factory({
      propsData: { types },
    });

    await flushPromises();

    expect(fetchPbehaviorTypesFieldList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: PAGINATION_LIMIT,
        search: '',
        with_hidden: true,
        types,
      },
    });
  });

  test('Hydrates the selected PBehavior type by id', async () => {
    const [type] = pbehaviorTypes;
    fetchPbehaviorTypesFieldList.mockImplementation(({ params }) =>
      Promise.resolve({
        data: params.ids ? [type] : [],
        meta: { page_count: 1 },
      }),
    );

    const wrapper = factory({
      propsData: {
        modelValue: type._id,
      },
    });

    await flushPromises();

    expect(fetchPbehaviorTypesFieldList).toHaveBeenCalledWith({
      params: {
        ids: [type._id],
        limit: 1,
        with_hidden: true,
      },
    });
    expect(wrapper.findComponent(lazySearchFieldStub).props('value')).toEqual(type);
  });

  test('Fetches the next PBehavior types page', async () => {
    fetchPbehaviorTypesFieldList.mockResolvedValue({
      data: pbehaviorTypes,
      meta: { page_count: 2 },
    });
    const wrapper = factory();

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('fetch:more');
    await flushPromises();

    expect(fetchPbehaviorTypesFieldList).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 2,
        limit: PAGINATION_LIMIT,
        with_hidden: true,
      }),
    });
  });

  test('Searches PBehavior types from the first page', async () => {
    vi.useFakeTimers();
    const wrapper = factory();

    await flushPromises();
    fetchPbehaviorTypesFieldList.mockClear();

    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:search', 'searched type');
    await vi.advanceTimersByTimeAsync(100);
    await flushPromises();

    expect(fetchPbehaviorTypesFieldList).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        search: 'searched type',
        with_hidden: true,
      }),
    });
  });

  test('Emits the selected PBehavior type id', async () => {
    fetchPbehaviorTypesFieldList.mockResolvedValueOnce({
      data: pbehaviorTypes,
      meta: { page_count: 1 },
    });
    const wrapper = factory();

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:value', pbehaviorTypes[0]);

    expect(wrapper).toEmitInput(pbehaviorTypes[0]._id);
  });

  test('Emits selected PBehavior type objects for a multiple object field', async () => {
    fetchPbehaviorTypesFieldList.mockResolvedValueOnce({
      data: pbehaviorTypes,
      meta: { page_count: 1 },
    });
    const wrapper = factory({
      propsData: {
        modelValue: [],
        multiple: true,
        returnObject: true,
      },
    });

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:value', pbehaviorTypes);

    expect(wrapper).toEmitInput(pbehaviorTypes);
  });
});
