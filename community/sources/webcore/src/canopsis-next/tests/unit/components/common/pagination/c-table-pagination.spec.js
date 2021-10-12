import Faker from 'faker';

import { mount, createVueInstance } from '@unit/utils/vue';

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
  'c-records-per-page': {
    template: `
      <input class="c-records-per-page" @input="$listeners.input(+$event.target.value)" />
    `,
  },
};

const factory = (options = {}) => mount(CTablePagination, {
  localVue,
  stubs,
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

    const recordsPerPage = wrapper.find('.c-records-per-page');

    recordsPerPage.setValue(rowsPerPage);

    const updateRowsPerPageEvents = wrapper.emitted('update:rows-per-page');

    expect(updateRowsPerPageEvents).toHaveLength(1);
    expect(updateRowsPerPageEvents[0]).toEqual([rowsPerPage]);
  });

  it('Renders `c-table-pagination` correctly', () => {
    const wrapper = mount(CTablePagination, {
      localVue,
      propsData: { page: 3, rowsPerPage: 10, totalItems: 100 },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
