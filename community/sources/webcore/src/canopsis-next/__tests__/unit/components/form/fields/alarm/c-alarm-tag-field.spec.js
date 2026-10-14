import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { PAGINATION_LIMIT } from '@/config';

import CAlarmTag from '@/components/forms/fields/alarm/c-alarm-tag-field.vue';

const { fetchAlarmTagsListWithoutStore, fetchAlarmTagsLabelsListWithoutStore } = vi.hoisted(() => ({
  fetchAlarmTagsListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
  fetchAlarmTagsLabelsListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: { page_count: 1 },
  }),
}));

vi.mock('@/hooks/store/modules/alarm-tag', () => ({
  useAlarmTag: () => ({
    fetchAlarmTagsListWithoutStore,
  }),
}));

vi.mock('@/hooks/store/modules/alarm-tag-label', () => ({
  useAlarmTagLabel: () => ({
    fetchAlarmTagsLabelsListWithoutStore,
  }),
}));

const stubs = {
  'c-chip': true,
  'c-lazy-search-field': createSelectInputStub('c-lazy-search-field'),
};

const selectLazySearchField = wrapper => wrapper.find('select.c-lazy-search-field');

describe('c-alarm-tag-field', () => {
  const items = [
    {
      _id: '1',
      value: 'Tag 1',
      color: '#444',
    },
    {
      _id: '2',
      value: 'Tag 2',
      color: '#222',
    },
    {
      _id: '3',
      value: 'Tag 3',
      color: '#000',
    },
  ];

  const factory = generateShallowRenderer(CAlarmTag, { stubs });

  afterEach(() => {
    fetchAlarmTagsListWithoutStore.mockClear();
    fetchAlarmTagsLabelsListWithoutStore.mockClear();
  });

  test('fetchAlarmTags was called', async () => {
    factory();

    await flushPromises();

    expect(fetchAlarmTagsListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        limit: PAGINATION_LIMIT,
      }),
    });
  });

  test('fetchAlarmTagsLabels was called', async () => {
    factory({
      propsData: {
        onlyLabels: true,
      },
    });

    await flushPromises();

    expect(fetchAlarmTagsLabelsListWithoutStore).toHaveBeenCalledWith({
      params: expect.objectContaining({
        page: 1,
        limit: PAGINATION_LIMIT,
      }),
    });
  });

  test('Value changed after trigger the input', () => {
    const wrapper = factory();

    selectLazySearchField(wrapper).triggerCustomEvent('update:value', items[0]);

    expect(wrapper).toEmitInput([items[0].value]);
  });
});
