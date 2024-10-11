import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import Faker from 'faker';

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
    const initialRowsPerPage = Faker.datatype.number();
    const wrapper = factory({
      data() {
        return {
          query: {
            rowsPerPage: initialRowsPerPage,
          },
        };
      },
      store,
      listeners,
    });

    await flushPromises();

    fetchPatternsList.mockClear();

    const patternsListNode = selectPattersListNode(wrapper);

    const rowsPerPage = Faker.datatype.number({ max: initialRowsPerPage });
    const page = Faker.datatype.number();

    patternsListNode.$emit('update:pagination', {
      rowsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchPatternsList).toBeCalledTimes(1);
    expect(fetchPatternsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: rowsPerPage,
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

  it('Renders `patterns` without patterns', () => {
    const wrapper = snapshotFactory({ store, listeners });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `patterns` with patterns', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
