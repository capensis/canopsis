import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { PAGINATION_LIMIT } from '@/config';

import CPbehaviorReasonField from '@/components/forms/fields/pbehavior/c-pbehavior-reason-field.vue';

const { fetchPbehaviorReasonsListWithoutStore } = vi.hoisted(() => ({
  fetchPbehaviorReasonsListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
}));

vi.mock('@/hooks/store/modules/pbehavior-reason', () => ({
  usePbehaviorReason: () => ({
    fetchPbehaviorReasonsListWithoutStore,
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

const pbehaviorReasons = [
  { _id: 'reason-1', name: 'Reason 1' },
  { _id: 'reason-2', name: 'Reason 2' },
];

describe('c-pbehavior-reason-field', () => {
  const factory = generateShallowRenderer(CPbehaviorReasonField, { stubs });

  afterEach(() => {
    vi.useRealTimers();
    fetchPbehaviorReasonsListWithoutStore.mockClear();
    fetchPbehaviorReasonsListWithoutStore.mockResolvedValue({
      data: [],
      meta: { page_count: 1 },
    });
  });

  test('Fetches the first PBehavior reasons page', async () => {
    factory();

    await flushPromises();

    expect(fetchPbehaviorReasonsListWithoutStore).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: PAGINATION_LIMIT,
        search: '',
        with_hidden: true,
      },
    });
  });

  test('Hydrates the selected PBehavior reason by id', async () => {
    const [reason] = pbehaviorReasons;
    fetchPbehaviorReasonsListWithoutStore.mockImplementation(({ params }) =>
      Promise.resolve({
        data: params.ids ? [reason] : [],
        meta: { page_count: 1 },
      }),
    );

    const wrapper = factory({
      propsData: {
        modelValue: reason._id,
      },
    });

    await flushPromises();

    expect(fetchPbehaviorReasonsListWithoutStore).toHaveBeenCalledWith({
      params: {
        ids: [reason._id],
        limit: 1,
        with_hidden: true,
      },
    });
    expect(wrapper.findComponent(lazySearchFieldStub).props('value')).toEqual(reason);
  });

  test('Fetches the next PBehavior reasons page', async () => {
    fetchPbehaviorReasonsListWithoutStore.mockResolvedValue({
      data: pbehaviorReasons,
      meta: { page_count: 2 },
    });
    const wrapper = factory();

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('fetch:more');
    await flushPromises();

    expect(fetchPbehaviorReasonsListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 2,
        limit: PAGINATION_LIMIT,
        with_hidden: true,
      }),
    });
  });

  test('Searches PBehavior reasons from the first page', async () => {
    vi.useFakeTimers();
    const wrapper = factory();

    await flushPromises();
    fetchPbehaviorReasonsListWithoutStore.mockClear();

    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:search', 'searched reason');
    await vi.advanceTimersByTimeAsync(100);
    await flushPromises();

    expect(fetchPbehaviorReasonsListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        search: 'searched reason',
        with_hidden: true,
      }),
    });
  });

  test('Emits the selected PBehavior reason id', async () => {
    fetchPbehaviorReasonsListWithoutStore.mockResolvedValueOnce({
      data: pbehaviorReasons,
      meta: { page_count: 1 },
    });
    const wrapper = factory();

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:value', pbehaviorReasons[0]);
    await flushPromises();

    expect(wrapper).toEmitInput(pbehaviorReasons[0]._id);
    expect(wrapper.emitted('update:selected-items').at(-1)).toEqual([pbehaviorReasons.slice(0, 1)]);
  });

  test('Emits the selected PBehavior reason object', async () => {
    fetchPbehaviorReasonsListWithoutStore.mockResolvedValueOnce({
      data: pbehaviorReasons,
      meta: { page_count: 1 },
    });
    const wrapper = factory({
      propsData: {
        returnObject: true,
      },
    });

    await flushPromises();
    wrapper.findComponent(lazySearchFieldStub).vm.$emit('update:value', pbehaviorReasons[0]);

    expect(wrapper).toEmitInput(pbehaviorReasons[0]);
  });
});
