import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { PAGINATION_LIMIT } from '@/config';

import CPbehaviorField from '@/components/forms/fields/pbehavior/c-pbehavior-field.vue';

const { fetchPbehaviorsListWithoutStore } = vi.hoisted(() => ({
  fetchPbehaviorsListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
}));

vi.mock('@/hooks/store/modules/pbehavior', () => ({
  usePbehavior: () => ({
    fetchPbehaviorsListWithoutStore,
  }),
}));

const lazySearchFieldBaseStub = createSelectInputStub('c-lazy-search-field');
const lazySearchFieldStub = {
  ...lazySearchFieldBaseStub,
  name: 'c-lazy-search-field',
  props: {
    ...lazySearchFieldBaseStub.props,
    itemTitle: [String, Function],
  },
  emits: [...lazySearchFieldBaseStub.emits, 'fetch', 'fetch:more', 'update:search'],
};

const stubs = {
  'c-lazy-search-field': lazySearchFieldStub,
};

const pbehaviors = [
  { _id: 'pbehavior-1', name: 'PBehavior 1' },
  { _id: 'pbehavior-2', name: 'PBehavior 2' },
];

describe('c-pbehavior-field', () => {
  const factory = generateShallowRenderer(CPbehaviorField, { stubs });

  afterEach(() => {
    vi.useRealTimers();
    fetchPbehaviorsListWithoutStore.mockClear();
    fetchPbehaviorsListWithoutStore.mockResolvedValue({
      data: [],
      meta: { page_count: 1 },
    });
  });

  test('Fetches the first PBehavior page', async () => {
    factory();

    await flushPromises();

    expect(fetchPbehaviorsListWithoutStore).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: PAGINATION_LIMIT,
        search: '',
      },
    });
  });

  test('Hydrates the selected PBehavior by id', async () => {
    const [pbehavior] = pbehaviors;
    fetchPbehaviorsListWithoutStore.mockImplementation(({ params }) =>
      Promise.resolve({
        data: params.ids ? [pbehavior] : [],
        meta: { page_count: 1 },
      }),
    );

    const wrapper = factory({
      propsData: {
        value: pbehavior._id,
      },
    });

    await flushPromises();

    expect(fetchPbehaviorsListWithoutStore).toHaveBeenCalledWith({
      params: {
        ids: [pbehavior._id],
        limit: 1,
      },
    });
    expect(wrapper.findComponent(lazySearchFieldStub).props('value')).toEqual(pbehavior);
  });

  test('Fetches the next PBehavior page', async () => {
    fetchPbehaviorsListWithoutStore.mockResolvedValue({
      data: pbehaviors,
      meta: { page_count: 2 },
    });
    const wrapper = factory();

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('fetch:more');
    await flushPromises();

    expect(fetchPbehaviorsListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 2,
        limit: PAGINATION_LIMIT,
      }),
    });
  });

  test('Searches PBehaviors from the first page', async () => {
    vi.useFakeTimers();
    const wrapper = factory();

    await flushPromises();
    fetchPbehaviorsListWithoutStore.mockClear();

    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:search', 'searched pbehavior');
    await vi.advanceTimersByTimeAsync(100);
    await flushPromises();

    expect(fetchPbehaviorsListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        search: 'searched pbehavior',
      }),
    });
  });

  test('Passes configurable item fields to the lazy search', () => {
    const wrapper = factory({
      propsData: {
        itemTitle: '_id',
        itemValue: 'name',
      },
    });
    const field = wrapper.findComponent(lazySearchFieldStub);

    expect(field.props('itemTitle')).toBe('_id');
    expect(field.props('itemValue')).toBe('name');
  });

  test('Emits the selected PBehavior id', async () => {
    fetchPbehaviorsListWithoutStore.mockResolvedValueOnce({
      data: pbehaviors,
      meta: { page_count: 1 },
    });
    const wrapper = factory();

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:value', pbehaviors[0]);

    expect(wrapper).toEmitInput(pbehaviors[0]._id);
    expect(wrapper.emitted('update:modelValue').at(-1)).toEqual([pbehaviors[0]._id]);
  });

  test('Emits selected PBehavior objects for a multiple object field', async () => {
    fetchPbehaviorsListWithoutStore.mockResolvedValueOnce({
      data: pbehaviors,
      meta: { page_count: 1 },
    });
    const wrapper = factory({
      propsData: {
        value: [],
        returnObject: true,
      },
    });

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:value', pbehaviors);

    expect(wrapper).toEmitInput(pbehaviors);
  });
});
