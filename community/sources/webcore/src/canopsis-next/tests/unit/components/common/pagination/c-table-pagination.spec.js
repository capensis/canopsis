import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';

const mockData = {
  page: Faker.datatype.number(),
  itemsPerPage: Faker.datatype.number(),
  falseTotalItems: 0,
  totalItems: Faker.datatype.number({ min: 1 }),
};

const stubs = {
  'c-pagination': {
    template: `
      <input class="c-pagination" @input="$listeners.input(+$event.target.value)" />
    `,
  },
  'c-items-per-page-field': {
    template: `
      <input class="c-items-per-page-field" @input="$listeners.input(+$event.target.value)" />
    `,
  },
};

const snapshotStubs = {
  'c-items-per-page-field': true,
  'c-pagination': true,
};

describe('c-table-pagination', () => {
  const factory = generateShallowRenderer(CTablePagination, { stubs });
  const snapshotFactory = generateRenderer(CTablePagination, { stubs: snapshotStubs });

  it('Pagination hidden without items', () => {
    const { falseTotalItems: totalItems } = mockData;
    const wrapper = factory({ propsData: { totalItems } });

    expect(wrapper.isVisible()).toBe(false);
  });

  it('Update pagination page', () => {
    const { page } = mockData;
    const wrapper = factory();

    wrapper.find('.c-pagination').setValue(page);

    expect(wrapper).toEmit('update:page', page);
  });

  it('Update pagination rows per page', () => {
    const { itemsPerPage } = mockData;
    const wrapper = factory();

    wrapper.find('.c-items-per-page-field').setValue(itemsPerPage);

    expect(wrapper).toEmit('update:items-per-page', itemsPerPage);
  });

  it('Renders `c-table-pagination` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: { page: 3, itemsPerPage: 10, totalItems: 100 },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
