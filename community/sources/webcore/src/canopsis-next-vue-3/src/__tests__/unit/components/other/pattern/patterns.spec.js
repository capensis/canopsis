import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { ref } from 'vue';

import Patterns from '@/components/other/pattern/patterns.vue';

vi.mock('@/hooks/store/modules/pattern', () => ({
  usePattern: () => globalThis.__TEST_PATTERN_HOOK__,
}));

const stubs = {
  'patterns-list': {
    name: 'PatternsList',
    emits: ['update:options', 'edit', 'remove', 'remove-selected'],
    template: '<div />',
  },
};

const selectPatternsList = wrapper => wrapper.findComponent({ name: 'PatternsList' });

describe('patterns', () => {
  const fetchPatternsList = vi.fn();

  const factory = generateShallowRenderer(Patterns, { stubs });

  beforeEach(() => {
    globalThis.__TEST_PATTERN_HOOK__ = {
      patterns: ref([]),
      patternsMeta: ref({ total_count: 0 }),
      patternsPending: ref(false),
      fetchPatternsList,
    };
  });

  afterEach(() => {
    fetchPatternsList.mockClear();
  });

  test('Filters fetched after mount', async () => {
    factory();

    await flushPromises();

    expect(fetchPatternsList).toHaveBeenCalledTimes(1);
    expect(fetchPatternsList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: 10,
        corporate: false,
      },
    });
  });

  test('Filters fetched after change query', async () => {
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    const page = Faker.datatype.number({ min: 2 });

    selectPatternsList(wrapper).vm.$emit('update:options', {
      page,
      itemsPerPage: 10,
    });

    await flushPromises();

    expect(fetchPatternsList).toHaveBeenCalledTimes(1);
    expect(fetchPatternsList).toHaveBeenCalledWith({
      params: {
        page,
        limit: 10,
        corporate: false,
      },
    });
  });

  test('Edit event emitted after trigger edit event on patterns list', async () => {
    const data = Faker.datatype.string();
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    selectPatternsList(wrapper).vm.$emit('edit', data);

    await flushPromises();

    expect(wrapper.emitted('edit')).toEqual([[data]]);
  });

  test('Remove selected event emitted after trigger remove selected event on patterns list', async () => {
    const data = [Faker.datatype.string()];
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    selectPatternsList(wrapper).vm.$emit('remove-selected', data);

    await flushPromises();

    expect(wrapper.emitted('remove-selected')).toEqual([[data]]);
  });

  test('Remove event emitted after trigger remove event on patterns list', async () => {
    const data = Faker.datatype.string();
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    selectPatternsList(wrapper).vm.$emit('remove', data);

    await flushPromises();

    expect(wrapper.emitted('remove')).toEqual([[data]]);
  });
});
