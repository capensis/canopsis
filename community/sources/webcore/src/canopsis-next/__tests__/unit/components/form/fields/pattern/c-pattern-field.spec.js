import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { MAX_LIMIT, PATTERN_CUSTOM_ITEM_VALUE, PATTERN_TYPES } from '@/constants';

import CPatternField from '@/components/forms/fields/pattern/c-pattern-field.vue';

const { fetchPatternsListWithoutStore } = vi.hoisted(() => ({
  fetchPatternsListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
  }),
}));

vi.mock('@/hooks/store/modules/pattern', () => ({
  usePattern: () => ({
    fetchPatternsListWithoutStore,
  }),
}));

const stubs = {
  'c-select': createSelectInputStub('c-select'),
};

const selectSelectField = wrapper => wrapper.find('select.c-select');

describe('c-pattern-field', () => {
  const factory = generateShallowRenderer(CPatternField, { stubs });

  afterEach(() => {
    fetchPatternsListWithoutStore.mockClear();
  });

  test('Patterns fetched after mount', async () => {
    factory();

    await flushPromises();

    expect(fetchPatternsListWithoutStore).toHaveBeenCalledWith({
      params: { limit: MAX_LIMIT },
    });
  });

  test('Pbehavior patterns fetched after mount', async () => {
    factory({
      propsData: {
        type: PATTERN_TYPES.pbehavior,
      },
    });

    await flushPromises();

    expect(fetchPatternsListWithoutStore).toHaveBeenCalledWith({
      params: { limit: MAX_LIMIT, type: PATTERN_TYPES.pbehavior },
    });
  });

  test('Value changed after trigger the select', () => {
    const wrapper = factory();

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', PATTERN_CUSTOM_ITEM_VALUE);

    expect(wrapper).toEmitInput(PATTERN_CUSTOM_ITEM_VALUE);
  });
});
