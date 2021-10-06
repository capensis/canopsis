import Faker from 'faker';

import { mount, createVueInstance } from '@/unit/utils/vue';

import CAdvancedSearchField from '@/components/forms/fields/c-advanced-search-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-search-field': {
    props: ['value'],
    template: `
      <div class='c-search-field'>
        <input
          :value="value"
          @input="$emit('input', $event.target.value)"
          @keydown="$emit('submit')"
        />
        <button class="c-search-field_submit" @click="$emit('submit')" />
        <button class="c-search-field_clear" @click="$emit('clear')" />
        <slot />
      </div>
    `,
  },
  'v-tooltip': {
    template: `
      <div class='v-tooltip'>
        <slot />
      </div>
    `,
  },
};

const factory = (options = {}) => mount(CAdvancedSearchField, {
  localVue,
  stubs,
  ...options,
});

describe('c-advanced-search-field', () => {
  it('Pass default search into props and it was applied to input field', () => {
    const search = Faker.lorem.words();

    const wrapper = factory({ propsData: { query: { search } } });
    const input = wrapper.find('input');

    expect(input.element.value).toBe(search);
  });

  it('Search hints was applied', () => {
    const tooltipText = Faker.lorem.words();
    const wrapper = factory({
      propsData: {
        query: {},
        tooltip: tooltipText,
      },
    });
    const tooltipElement = wrapper.find('.v-tooltip');

    expect(tooltipElement.html()).toContain(tooltipText);
  });

  it('Search hints as html was applied', () => {
    const tooltipContentText = Faker.lorem.words();
    const tooltip = `
      <span class="tooltip-content-span">${tooltipContentText}</span>
    `;
    const wrapper = factory({
      propsData: {
        query: {},
        tooltip,
      },
    });
    const tooltipSpanElement = wrapper.find('.tooltip-content-span');

    expect(tooltipSpanElement.html()).toContain(tooltipContentText);
  });

  it('Submit search button', () => {
    const search = Faker.lorem.words();

    const wrapper = factory({ propsData: { query: { search } } });
    const input = wrapper.find('input');

    input.trigger('keydown');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
  });

  it('Submit search with search by column without columns', () => {
    const column = {
      text: Faker.lorem.word(),
      value: Faker.lorem.word(),
    };
    const searchNameValue = Faker.lorem.words();
    const search = `- ${column.text} = '${searchNameValue}'`;

    const wrapper = factory({
      propsData: {
        query: { search },
      },
    });
    const submitButton = wrapper.find('.c-search-field_submit');

    submitButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
      search: `${column.text} = '${searchNameValue}'`,
    }]);
  });

  it('Submit search with search by column', () => {
    const column = {
      text: Faker.lorem.word(),
      value: Faker.lorem.word(),
    };
    const searchNameValue = Faker.lorem.words();
    const search = `- ${column.text} = '${searchNameValue}'`;

    const wrapper = factory({
      propsData: {
        query: { search },
        columns: [column],
      },
    });
    const submitButton = wrapper.find('.c-search-field_submit');

    submitButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
      search: `${column.value} = '${searchNameValue}'`,
    }]);
  });

  it('Clear search by click on button', () => {
    const search = Faker.lorem.words();
    const wrapper = factory({
      propsData: {
        query: { search },
      },
    });
    const clearButton = wrapper.find('.c-search-field_clear');

    clearButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{}]);
  });

  it('Submit search with custom field name', () => {
    const search = Faker.lorem.words();
    const field = Faker.lorem.word();
    const wrapper = factory({
      propsData: {
        query: { [field]: '' },
        field,
      },
    });
    const submitButton = wrapper.find('.c-search-field_submit');
    const input = wrapper.find('input');

    input.setValue(search);

    submitButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
      [field]: search,
    }]);
  });

  it('Clear search with custom field name', () => {
    const search = Faker.lorem.words();
    const field = Faker.lorem.word();
    const wrapper = factory({
      propsData: {
        query: { [field]: search },
        field,
      },
    });
    const clearButton = wrapper.find('.c-search-field_clear');

    clearButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{}]);
  });

  it('Renders `c-advanced-search-field` correctly', () => {
    const tooltip = `
      <span>Tooltip content</span>
    `;
    const wrapper = mount(CAdvancedSearchField, {
      localVue,
      stubs: ['c-search-field'],
      propsData: {
        query: { search: 'c-advanced-search-field' },
        tooltip,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
