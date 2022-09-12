import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';

const localVue = createVueInstance();

const mockData = {
  page: Faker.datatype.number(),
  rowsPerPage: Faker.datatype.number(),
  falseTotalItems: 0,
  totalItems: Faker.datatype.number({ min: 1 }),
};

const stubs = {
  'c-pagination': {
    template: `
      <input class="c-pagination" @input="$listeners.input(+$event.target.value)" />
    `,
  },
  'c-records-per-page-field': {
    template: `
      <input class="c-records-per-page-field" @input="$listeners.input(+$event.target.value)" />
    `,
  },
};

const snapshotStubs = {
  'c-records-per-page-field': true,
  'c-pagination': true,
};

const factory = (options = {}) => shallowMount(CTablePagination, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CTablePagination, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('c-table-pagination', () => {
  it('Pagination hidden without items', () => {
    const { falseTotalItems: totalItems } = mockData;
    const wrapper = factory({ propsData: { totalItems } });

    expect(wrapper.isVisible()).toBe(false);
  });

  it('Update pagination page', () => {
    const { page } = mockData;
    const wrapper = factory();

    const pagination = wrapper.find('.c-pagination');

    pagination.setValue(page);

    const updatePageEvents = wrapper.emitted('update:page');

    expect(updatePageEvents).toHaveLength(1);
    expect(updatePageEvents[0]).toEqual([page]);
  });

  it('Update pagination rows per page', () => {
    const { rowsPerPage } = mockData;
    const wrapper = factory();

    const recordsPerPage = wrapper.find('.c-records-per-page-field');

    recordsPerPage.setValue(rowsPerPage);

    const updateRowsPerPageEvents = wrapper.emitted('update:rows-per-page');

    expect(updateRowsPerPageEvents).toHaveLength(1);
    expect(updateRowsPerPageEvents[0]).toEqual([rowsPerPage]);
  });

  it('Renders `c-table-pagination` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: { page: 3, rowsPerPage: 10, totalItems: 100 },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
