import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';

const mockData = {
  page: Faker.datatype.number(),
  falseTotalItems: 0,
  totalItems: Faker.datatype.number({ min: 1 }),
};

const stubs = {
  'c-pagination': {
    template: `
      <input class="c-pagination" @input="$emit('input', +$event.target.value)" />
    `,
  },
  'c-items-per-page-field': {
    template: `
      <input class="c-items-per-page-field" @input="$emit('input', +$event.target.value)" />
    `,
  },
};

describe('c-table-pagination', () => {
  const factory = generateShallowRenderer(CTablePagination, { stubs });

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

  it('Update pagination items per page (without page changing)', () => {
    const itemsPerPage = 10;
    const wrapper = factory({
      propsData: {
        page: 1,
        itemsPerPage: 5,
      },
    });

    wrapper.find('.c-items-per-page-field').setValue(itemsPerPage);

    expect(wrapper).toEmit('update:items-per-page', itemsPerPage);
  });

  it('Update pagination items per page (with page changing)', () => {
    const itemsPerPage = 20;
    const wrapper = factory({
      propsData: {
        page: 2,
        itemsPerPage: 5,
      },
    });

    wrapper.find('.c-items-per-page-field').setValue(itemsPerPage);

    expect(wrapper).toEmitInput({ page: 1, itemsPerPage });
  });
});
