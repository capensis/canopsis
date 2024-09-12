import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import VDataTable from '@/plugins/vuetify/components/v-data-table/v-data-table.vue';

const stubs = {};

const selectTableHeader = wrapper => wrapper.findAll('th');

describe('v-data-table', () => {
  const headers = [
    { value: 'column', sortable: true },
    { value: 'column2' },
    { value: 'column3', sortable: true },
  ];
  const [sortableHeader, , sortableHeaderTwo] = headers;
  const rowsPerPage = 5;
  const page = 1;
  const totalItems = 0;

  const factory = generateShallowRenderer(VDataTable, { stubs });
  const snapshotFactory = generateRenderer(VDataTable, { stubs });

  test('Column sorted by DESC after click on the header', () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual({
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
    });

    const tableHeader = selectTableHeader(wrapper).at(0);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      page,
      rowsPerPage,
      totalItems,

      descending: true,
      sortBy: sortableHeader.value,
    });
  });

  test('Other column sorted by ASC after click on the header', () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual({
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
    });

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeaderTwo.value,
    });
  });

  test('Column sorted by ASC after trigger keydown on the header', () => {
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual({
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
    });

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('keydown', { keyCode: 32 });

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeaderTwo.value,
    });
  });

  test('Column sort reset after click on the header', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: true,
      sortBy: sortableHeaderTwo.value,
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        pagination: {
          page,
          rowsPerPage,
          totalItems,

          descending: true,
          sortBy: sortableHeaderTwo.value,
        },
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      descending: null,
      sortBy: null,
    });
  });

  test('Column sorted by DESC after click on the header with must sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: true,
      sortBy: sortableHeaderTwo.value,
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        mustSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,
      descending: false,
    });
  });

  test('First column sorted by DESC after click on the header with multi sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: true,
      sortBy: sortableHeader.value,
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(0);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
      ],
    });
  });

  test('First column sorted by ASC after click on the header with multi sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(0);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      multiSortBy: [
        {
          descending: true,
          sortBy: sortableHeader.value,
        },
      ],
    });
  });

  test('Second column sorted by DESC after click on the header with multi sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
        {
          descending: false,
          sortBy: sortableHeaderTwo.value,
        },
      ],
    });
  });

  test('Second column sorted by ASC after click on the header with multi sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
        {
          descending: false,
          sortBy: sortableHeaderTwo.value,
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
        {
          descending: true,
          sortBy: sortableHeaderTwo.value,
        },
      ],
    });
  });

  test('Second column sort reset after click on the header with multi sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
        {
          descending: true,
          sortBy: sortableHeaderTwo.value,
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(2);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      multiSortBy: [
        {
          descending: false,
          sortBy: sortableHeader.value,
        },
      ],
    });
  });

  test('Column sort reset after click on the header with multi sort', () => {
    const pagination = {
      page,
      rowsPerPage,
      totalItems,

      descending: false,
      sortBy: sortableHeader.value,
      multiSortBy: [
        {
          descending: true,
          sortBy: sortableHeader.value,
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        pagination,
      },
    });

    const updatePaginationEventsAfterMount = wrapper.emitted('update:pagination');
    const [eventDataAfterMount] = updatePaginationEventsAfterMount[0];
    expect(eventDataAfterMount).toEqual(pagination);

    const tableHeader = selectTableHeader(wrapper).at(0);

    tableHeader.trigger('click');

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    const [eventData] = updatePaginationEvents[1];
    expect(eventData).toEqual({
      ...pagination,

      multiSortBy: [],
    });
  });

  test('Renders `v-data-table` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `v-data-table` with hidden headers', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
        hideHeaders: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `v-data-table` with custom headers', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
      },
      scopedSlots: {
        headers(props) {
          return this.$createElement(
            'tr',
            `Custom headers${JSON.stringify(props)}`,
          );
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `v-data-table` with select all', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
        selectAll: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `v-data-table` with items and headers', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [{ property1: 'Property 1', property2: 'Property 2' }],
        headers: [
          { value: 'property1', text: 'Header 1' },
          { value: 'property2', text: 'Header 2', sortable: true },
        ],
      },
      scopedSlots: {
        items(props) {
          return this.$createElement('tr', `Row${JSON.stringify(props)}`);
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `v-data-table` with items, headers and ellipsis headers', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [{ property1: 'Property 1', property2: 'Property 2' }],
        headers: [
          { value: 'property1', text: 'Header 1' },
          { value: 'property2', text: 'Header 2', sortable: true },
        ],
        ellipsisHeaders: true,
      },
      scopedSlots: {
        items(props) {
          return this.$createElement('tr', `Row${JSON.stringify(props)}`);
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
