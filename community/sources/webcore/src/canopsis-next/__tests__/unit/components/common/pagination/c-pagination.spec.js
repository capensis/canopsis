import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CPagination from '@/components/common/pagination/c-pagination.vue';

const stubs = {
  'v-pagination': {
    emits: ['update:modelValue'],
    template: `
      <input class="v-pagination" @input="$emit('update:modelValue', +$event.target.value)" />
    `,
  },
};

describe('c-pagination', () => {
  const factory = generateShallowRenderer(CPagination, { stubs });

  it('Pagination hidden without total', () => {
    const wrapper = factory({ propsData: { total: 0 } });

    expect(wrapper.isVisible()).toBe(false);
  });

  it('Pagination on the top. Check prev page button.', () => {
    const page = 2;
    const wrapper = factory({ propsData: { page, limit: 5, total: 10, type: 'top' } });

    wrapper.findAll('button').at(0).trigger('click');

    expect(wrapper).toEmitInput(page - 1);
  });

  it('Pagination on the top. Check prev page button is disabled, when first page.', () => {
    const wrapper = factory({ propsData: { total: 10, page: 1, type: 'top' } });

    const prevButton = wrapper.findAll('button').at(0);

    expect(prevButton.attributes('disabled')).toBeDefined();
  });

  it('Pagination on the top. Check next page button.', () => {
    const page = 1;
    const wrapper = factory({ propsData: { page, total: 10, limit: 5, type: 'top' } });

    wrapper.findAll('button').at(1).trigger('click');

    expect(wrapper).toEmitInput(page + 1);
  });

  it('Pagination on the top. Check next page button is disabled, when last page.', () => {
    const page = 2;
    const wrapper = factory({ propsData: { page, total: 10, limit: 5, type: 'top' } });

    const nextButton = wrapper.findAll('button').at(1);

    expect(nextButton.attributes('disabled')).toBeDefined();
  });

  it('Pagination on the bottom. Check next page button is disabled, when last page.', () => {
    const page = Faker.datatype.number();
    const wrapper = factory({ propsData: { total: 1 } });

    wrapper.find('.v-pagination').setValue(page);

    expect(wrapper).toEmitInput(page);
  });
});
