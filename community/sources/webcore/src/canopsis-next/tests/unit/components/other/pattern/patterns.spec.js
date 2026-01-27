import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import Patterns from '@/components/other/pattern/patterns.vue';

const stubs = {
  'patterns-list': true,
};

const selectPattersListNode = wrapper => wrapper.vm.$children[0];

describe('patterns', () => {
  const patterns = [
    { _id: 'id', title: 'title' },
  ];

  const fetchPatternsList = jest.fn();

  const patternModule = {
    name: 'pattern',
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

  const factory = generateShallowRenderer(Patterns, { stubs });
  const snapshotFactory = generateRenderer(Patterns, { stubs });

  afterEach(() => {
    fetchPatternsList.mockClear();
  });

  test('Filters fetched after mount', async () => {
    factory({ store, listeners });

    await flushPromises();

    expect(fetchPatternsList).toHaveBeenCalledTimes(1);
    expect(fetchPatternsList).toHaveBeenCalledWith(
      expect.any(Object),
      {
        params: {
          page: 1,
          limit: 10,
          corporate: false,
        },
      },
    );
  });

  test('Filters fetched after change query', async () => {
    const wrapper = factory({
      store,
      listeners,
    });

    await flushPromises();

    fetchPatternsList.mockClear();

    const patternsListNode = selectPattersListNode(wrapper);
    const page = Faker.datatype.number({ min: 2 });

    patternsListNode.$emit('update:options', {
      page,
      itemsPerPage: 10,
    });

    await flushPromises();

    expect(fetchPatternsList).toHaveBeenCalledTimes(1);
    expect(fetchPatternsList).toHaveBeenCalledWith(
      expect.any(Object),
      {
        params: {
          page,
          limit: 10,
          corporate: false,
        },
      },
    );
  });

  test('Edit event emitted after trigger edit event on patterns list', async () => {
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

    expect(edit).toHaveBeenCalledTimes(1);
    expect(edit).toHaveBeenCalledWith(data);
  });

  test('Remove selected event emitted after trigger remove selected event on patterns list', async () => {
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

    expect(removeSelected).toHaveBeenCalledTimes(1);
    expect(removeSelected).toHaveBeenCalledWith(data);
  });

  test('Remove event emitted after trigger remove event on patterns list', async () => {
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

    expect(remove).toHaveBeenCalledTimes(1);
    expect(remove).toHaveBeenCalledWith(data);
  });

  test('Renders `patterns` without patterns', () => {
    const wrapper = snapshotFactory({ store, listeners });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `patterns` with patterns', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          ...patternModule,
          getters: {
            items: patterns,
            meta: { total_count: patterns.length },
            pending: false,
          },
        },
      ]),
      listeners,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
