import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { ref } from 'vue';

import CorporatePatterns from '@/components/other/pattern/corporate-patterns.vue';

vi.mock('@/hooks/store/modules/pattern-corporate', () => ({
  useCorporatePattern: () => globalThis.__TEST_CORPORATE_PATTERN_HOOK__,
}));

const stubs = {
  'patterns-list': {
    name: 'PatternsList',
    emits: ['update:options', 'edit', 'remove', 'remove-selected'],
    template: '<div />',
  },
};

const selectPatternsList = wrapper => wrapper.findComponent({ name: 'PatternsList' });

describe('corporate-patterns', () => {
  const fetchPatternsList = vi.fn();

  const factory = generateShallowRenderer(CorporatePatterns, { stubs });

  beforeEach(() => {
    globalThis.__TEST_CORPORATE_PATTERN_HOOK__ = {
      corporatePatterns: ref([]),
      corporatePatternsMeta: ref({ total_count: 0 }),
      corporatePatternsPending: ref(false),
      fetchCorporatePatternsList: fetchPatternsList,
    };
  });

  afterEach(() => {
    fetchPatternsList.mockClear();
  });

  it('Filters fetched after mount', async () => {
    factory();

    await flushPromises();

    expect(fetchPatternsList).toBeCalledTimes(1);
    expect(fetchPatternsList).toHaveBeenCalledWith({
      params: {
        limit: 10,
        page: 1,
      },
    });
  });

  it('Filters fetched after change query', async () => {
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    const itemsPerPage = 5;
    const page = Faker.datatype.number({ min: 2 });

    selectPatternsList(wrapper).vm.$emit('update:options', {
      itemsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchPatternsList).toBeCalledTimes(1);
    expect(fetchPatternsList).toHaveBeenCalledWith({
      params: {
        limit: itemsPerPage,
        page,
      },
    });
  });

  it('Edit event emitted after trigger edit event on patterns list', async () => {
    const data = Faker.datatype.string();
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    selectPatternsList(wrapper).vm.$emit('edit', data);

    await flushPromises();

    expect(wrapper.emitted('edit')).toEqual([[data]]);
  });

  it('Remove selected event emitted after trigger remove selected event on patterns list', async () => {
    const data = [Faker.datatype.string()];
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    selectPatternsList(wrapper).vm.$emit('remove-selected', data);

    await flushPromises();

    expect(wrapper.emitted('remove-selected')).toEqual([[data]]);
  });

  it('Remove event emitted after trigger remove event on patterns list', async () => {
    const data = Faker.datatype.string();
    const wrapper = factory();

    await flushPromises();

    fetchPatternsList.mockClear();

    selectPatternsList(wrapper).vm.$emit('remove', data);

    await flushPromises();

    expect(wrapper.emitted('remove')).toEqual([[data]]);
  });
});
