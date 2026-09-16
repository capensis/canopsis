import { generateRenderer } from '@unit/utils/vue';

import VDataTable from '@/plugins/vuetify/components/v-data-table/v-data-table.vue';

const selectTableHeader = wrapper => wrapper.findAll('th');
const getOptionsEvents = wrapper => wrapper.emitted('update:options').map(([options]) => options);
const getSortEvents = wrapper => wrapper.emitted('update:sortBy').map(([sortBy]) => sortBy);

describe('v-data-table', () => {
  const headers = [
    { key: 'column', title: 'column', sortable: true },
    { key: 'column2', title: 'column2', sortable: false },
    { key: 'column3', title: 'column3', sortable: true },
  ];
  const [sortableHeader, , sortableHeaderTwo] = headers;
  const itemsPerPage = 5;
  const page = 1;

  const factory = generateRenderer(VDataTable, { stubs: {} });

  test('Default options emitted after mount', () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    expect(getOptionsEvents(wrapper)).toContainEqual({
      page,
      itemsPerPage: 10,
      sortBy: [],
      groupBy: [],
    });
  });

  test('Column sorted by ASC after click on the header', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    await selectTableHeader(wrapper).at(0).trigger('click');

    const sortBy = [{ key: sortableHeader.key, order: 'asc' }];

    expect(getSortEvents(wrapper)).toContainEqual(sortBy);
    expect(getOptionsEvents(wrapper)).toContainEqual({
      page,
      itemsPerPage: 10,
      sortBy,
      groupBy: [],
    });
  });

  test('Other column sorted by ASC after click on the header', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(getSortEvents(wrapper)).toContainEqual([{ key: sortableHeaderTwo.key, order: 'asc' }]);
  });

  test('Column sort changed to DESC after second click on the header', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        sortBy: [{ key: sortableHeaderTwo.key, order: 'asc' }],
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(getSortEvents(wrapper)).toContainEqual([{ key: sortableHeaderTwo.key, order: 'desc' }]);
  });

  test('Column sort reset after third click on the header', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        sortBy: [{ key: sortableHeaderTwo.key, order: 'desc' }],
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(getSortEvents(wrapper)).toContainEqual([]);
  });

  test('Column keeps sorting when must sort is enabled', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        mustSort: true,
        sortBy: [{ key: sortableHeaderTwo.key, order: 'desc' }],
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(getSortEvents(wrapper)).not.toContainEqual([]);
  });

  test('First column added to sort list with multi sort', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        sortBy: [{ key: sortableHeaderTwo.key, order: 'desc' }],
      },
    });

    await selectTableHeader(wrapper).at(0).trigger('click');

    expect(getSortEvents(wrapper)).toContainEqual([
      { key: sortableHeaderTwo.key, order: 'desc' },
      { key: sortableHeader.key, order: 'asc' },
    ]);
  });

  test('Existing column sort changed with multi sort', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        sortBy: [
          { key: sortableHeader.key, order: 'asc' },
          { key: sortableHeaderTwo.key, order: 'asc' },
        ],
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(getSortEvents(wrapper)).toContainEqual([
      { key: sortableHeader.key, order: 'asc' },
      { key: sortableHeaderTwo.key, order: 'desc' },
    ]);
  });

  test('Existing column sort reset with multi sort', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        sortBy: [
          { key: sortableHeader.key, order: 'asc' },
          { key: sortableHeaderTwo.key, order: 'desc' },
        ],
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(getSortEvents(wrapper)).toContainEqual([{ key: sortableHeader.key, order: 'asc' }]);
  });

  test('Non sortable column does not change sort', async () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    await selectTableHeader(wrapper).at(1).trigger('click');

    expect(wrapper.emitted('update:sortBy')).toBeUndefined();
  });

  test('Custom pagination options emitted after mount', () => {
    const sortBy = [{ key: sortableHeader.key, order: 'asc' }];
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        page: 2,
        itemsPerPage,
        sortBy,
      },
    });

    expect(getOptionsEvents(wrapper)).toContainEqual({
      page: 2,
      itemsPerPage,
      sortBy,
      groupBy: [],
    });
  });
});
