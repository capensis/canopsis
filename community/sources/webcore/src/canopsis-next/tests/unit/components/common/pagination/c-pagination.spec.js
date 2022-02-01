import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CPagination from '@/components/common/pagination/c-pagination.vue';

const localVue = createVueInstance();

const stubs = {
  'v-pagination': {
    template: `
      <input class="v-pagination" @input="$listeners.input(+$event.target.value)" />
    `,
  },
};

const factory = (options = {}) => shallowMount(CPagination, {
  localVue,
  stubs,
  ...options,
});

describe('c-pagination', () => {
  it('Pagination hidden without total', () => {
    const wrapper = factory({ propsData: { total: 0 } });

    expect(wrapper.isVisible()).toBe(false);
  });

  it('Pagination on the top. Check prev page button.', () => {
    const page = 2;
    const wrapper = factory({ propsData: { page, limit: 5, total: 10, type: 'top' } });

    const prevButton = wrapper.findAll('button').at(0);

    prevButton.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([page - 1]);
  });

  it('Pagination on the top. Check prev page button is disabled, when first page.', () => {
    const wrapper = factory({ propsData: { total: 10, page: 1, type: 'top' } });

    const prevButton = wrapper.findAll('button').at(0);

    expect(prevButton.attributes('disabled')).toBeTruthy();
  });

  it('Pagination on the top. Check next page button.', () => {
    const page = 1;
    const wrapper = factory({ propsData: { page, total: 10, limit: 5, type: 'top' } });

    const nextButton = wrapper.findAll('button').at(1);

    nextButton.trigger('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([page + 1]);
  });

  it('Pagination on the top. Check next page button is disabled, when last page.', () => {
    const page = 2;
    const wrapper = factory({ propsData: { page, total: 10, limit: 5, type: 'top' } });

    const nextButton = wrapper.findAll('button').at(1);

    expect(nextButton.attributes('disabled')).toBeTruthy();
  });

  it('Pagination on the bottom. Check next page button is disabled, when last page.', () => {
    const page = Faker.datatype.number();
    const wrapper = factory({ propsData: { total: 1 } });

    const pagination = wrapper.find('.v-pagination');

    pagination.setValue(page);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([page]);
  });

  it('Renders `c-pagination` with default props correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { total: 1 },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` with default props correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { total: 1 },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` on the top with disabled prev button correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { page: 1, total: 2, limit: 1, type: 'top' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` on the top with disabled next button correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { page: 2, total: 2, limit: 1, type: 'top' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` on the top with disabled buttons correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { page: 1, total: 1, limit: 1, type: 'top' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` with default props on the top correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { total: 1, type: 'top' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` on the bottom correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { page: 3, total: 100 },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-pagination` on the top correctly', () => {
    const wrapper = mount(CPagination, {
      localVue,
      propsData: { page: 3, total: 100, type: 'top' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
