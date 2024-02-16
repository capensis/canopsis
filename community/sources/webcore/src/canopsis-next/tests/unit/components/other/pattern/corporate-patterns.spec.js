import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import CorporatePatterns from '@/components/other/pattern/corporate-patterns.vue';

const stubs = {
  'patterns-list': true,
};

const selectPattersListNode = wrapper => wrapper.vm.$children[0];

describe('corporate-patterns', () => {
  const corporatePatterns = [
    { _id: 'id', title: 'title' },
  ];

  const fetchPatternsList = jest.fn();

  const patternModule = {
    name: 'pattern/corporate',
    getters: {
      pending: false,
      items: [],
      meta: {
        total_count: 0,
      },
    },
    actions: {
      fetchList: fetchPatternsList,
    },
  };
  const store = createMockedStoreModules([
    patternModule,
  ]);
  const edit = jest.fn();
  const remove = jest.fn();
  const removeSelected = jest.fn();
  const listeners = {
    edit,
    remove,
    'remove-selected': removeSelected,
  };

  const factory = generateShallowRenderer(CorporatePatterns, { stubs });
  const snapshotFactory = generateRenderer(CorporatePatterns, { stubs });

  afterEach(() => {
    fetchPatternsList.mockClear();
  });

  it('Filters fetched after mount', async () => {
    factory({ store, listeners });

    await flushPromises();

    expect(fetchPatternsList).toBeCalledTimes(1);
    expect(fetchPatternsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: 10,
          page: 1,
        },
      },
      undefined,
    );
  });

  it('Filters fetched after change query', async () => {
    const initialItemsPerPage = Faker.datatype.number();
    const wrapper = factory({
      data() {
        return {
          query: {
            itemsPerPage: initialItemsPerPage,
          },
        };
      },
      store,
      listeners,
    });

    await flushPromises();

    fetchPatternsList.mockClear();

    const patternsListNode = selectPattersListNode(wrapper);

    const itemsPerPage = Faker.datatype.number({ max: initialItemsPerPage });
    const page = Faker.datatype.number();

    patternsListNode.$emit('update:options', {
      itemsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchPatternsList).toBeCalledTimes(1);
    expect(fetchPatternsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: itemsPerPage,
          page,
        },
      },
      undefined,
    );
  });

  it('Edit event emitted after trigger edit event on patterns list', async () => {
    const data = Faker.datatype.string();
    const wrapper = factory({
      store,
      listeners,
    });

    await flushPromises();

    fetchPatternsList.mockClear();

    const patternsListNode = selectPattersListNode(wrapper);

    patternsListNode.$emit('edit', data);

    await flushPromises();

    expect(edit).toBeCalledTimes(1);
    expect(edit).toBeCalledWith(data);
  });

  it('Remove selected event emitted after trigger remove selected event on patterns list', async () => {
    const data = [Faker.datatype.string()];
    const wrapper = factory({
      store,
      listeners,
    });

    await flushPromises();

    fetchPatternsList.mockClear();

    const patternsListNode = selectPattersListNode(wrapper);

    patternsListNode.$emit('remove-selected', data);

    await flushPromises();

    expect(removeSelected).toBeCalledTimes(1);
    expect(removeSelected).toBeCalledWith(data);
  });

  it('Remove event emitted after trigger remove event on patterns list', async () => {
    const data = Faker.datatype.string();
    const wrapper = factory({
      store,
      listeners,
    });

    await flushPromises();

    fetchPatternsList.mockClear();

    const patternsListNode = selectPattersListNode(wrapper);

    patternsListNode.$emit('remove', data);

    await flushPromises();

    expect(remove).toBeCalledTimes(1);
    expect(remove).toBeCalledWith(data);
  });

  it('Renders `corporate-patterns` without patterns', () => {
    const wrapper = snapshotFactory({ store, listeners });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `corporate-patterns` with patterns', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          ...patternModule,
          getters: {
            items: corporatePatterns,
            meta: { total_count: corporatePatterns.length },
            pending: false,
          },
        },
      ]),
      listeners,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
