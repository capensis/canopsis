import { generateRenderer } from '@unit/utils/vue';

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
  const itemsPerPage = 5;
  const page = 1;
  const totalItems = 0;

  const snapshotFactory = generateRenderer(VDataTable, { stubs });

  it('Column sorted by DESC after click on the header', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        options: {},
      },
    });

    expect(wrapper).toEmit('update:options', {
      page,
      itemsPerPage: 10,

      multiSort: false,
      mustSort: false,
      sortDesc: [false],
      sortBy: [],

      groupBy: [],
      groupDesc: [],
    });

    await selectTableHeader(wrapper).at(0).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      expect.any(Object),
      {
        page,
        itemsPerPage: 10,

        multiSort: false,
        mustSort: false,
        sortDesc: [false],
        sortBy: [sortableHeader.value],

        groupBy: [],
        groupDesc: [],
      },
    );
  });

  it('Other column sorted by ASC after click on the header', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
      },
    });

    expect(wrapper).toEmit('update:options', {
      page,
      itemsPerPage: 10,

      multiSort: false,
      mustSort: false,
      sortDesc: [false],
      sortBy: [],

      groupBy: [],
      groupDesc: [],
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      expect.any(Object),
      {
        page,
        itemsPerPage: 10,

        multiSort: false,
        mustSort: false,
        sortDesc: [false],
        sortBy: [sortableHeaderTwo.value],

        groupBy: [],
        groupDesc: [],
      },
    );
  });

  it('Column sort reset after click on the header', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        options: {
          page,
          itemsPerPage,
          totalItems,

          sortDesc: [false],
          sortBy: [sortableHeaderTwo.value],
        },
      },
    });

    const newOptions = {
      page,
      itemsPerPage,
      totalItems,

      sortDesc: [false],
      sortBy: [sortableHeaderTwo.value],
      multiSort: false,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };

    await selectTableHeader(wrapper).at(2).trigger('click');
    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      newOptions,
      {
        ...newOptions,

        sortDesc: [true],
      },
      {
        ...newOptions,

        sortBy: [],
        sortDesc: [],
      },
    );
  });

  it('Column sorted by DESC after click on the header with must sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortDesc: [true],
      sortBy: [sortableHeaderTwo.value],
      multiSort: false,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        mustSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortBy: [],
        sortDesc: [],
      },
    );
  });

  it('First column sorted by DESC after click on the header with multi sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortDesc: [true],
      sortBy: [sortableHeaderTwo.value],
      multiSort: true,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(0).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortBy: [
          sortableHeaderTwo.value,
          sortableHeader.value,
        ],
        sortDesc: [
          true,
          false,
        ],
      },
    );
  });

  it('First column sorted by ASC after click on the header with multi sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortDesc: [false],
      sortBy: [sortableHeader.value],
      multiSort: true,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(0).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortDesc: [
          true,
        ],
      },
    );
  });

  it('Second column sorted by DESC after click on the header with multi sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortDesc: [false],
      sortBy: [sortableHeader.value],
      multiSort: true,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortBy: [sortableHeader.value, sortableHeaderTwo.value],
        sortDesc: [false, false],
      },
    );
  });

  it('Second column sorted by ASC after click on the header with multi sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortBy: [
        sortableHeader.value,
        sortableHeaderTwo.value,
      ],
      sortDesc: [false, false],
      multiSort: true,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortBy: [
          sortableHeader.value,
          sortableHeaderTwo.value,
        ],
        sortDesc: [false, true],
      },
    );
  });

  it('Second column sort reset after click on the header with multi sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortBy: [
        sortableHeader.value,
        sortableHeaderTwo.value,
      ],
      sortDesc: [false, true],
      multiSort: true,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(2).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortBy: [
          sortableHeader.value,
        ],
        sortDesc: [false],
      },
    );
  });

  it('Column sort reset after click on the header with multi sort', async () => {
    const pagination = {
      page,
      itemsPerPage,
      totalItems,

      sortBy: [
        sortableHeader.value,
      ],
      sortDesc: [true],
      multiSort: true,
      mustSort: false,

      groupBy: [],
      groupDesc: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: [],
        multiSort: true,
        options: pagination,
      },
    });

    await selectTableHeader(wrapper).at(0).trigger('click');

    expect(wrapper).toEmit(
      'update:options',
      pagination,
      {
        ...pagination,

        sortBy: [],
        sortDesc: [],
      },
    );
  });

  it('Renders `v-data-table` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `v-data-table` with hidden headers', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
        hideHeaders: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `v-data-table` with custom headers', () => {
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

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `v-data-table` with select all', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
        selectAll: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `v-data-table` with data and headers', () => {
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

    expect(wrapper).toMatchSnapshot();
  });
});
