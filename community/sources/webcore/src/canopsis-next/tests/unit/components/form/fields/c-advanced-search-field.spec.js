import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createActivatorElementStub } from '@unit/stubs/vuetify';

import CAdvancedSearchField from '@/components/forms/fields/c-advanced-search-field.vue';

const localVue = createVueInstance();

const mockData = {
  search: Faker.lorem.words(),
  tooltipText: Faker.lorem.words(),
  column: {
    text: 'Column label',
    value: 'column.value',
  },
  field: Faker.lorem.word(),
};

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
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

const factory = (options = {}) => shallowMount(CAdvancedSearchField, {
  localVue,
  stubs,
  ...options,
});

describe('c-advanced-search-field', () => {
  it('Pass default search into props and it was applied to input field', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { query: { search } } });
    const input = wrapper.find('input');

    expect(input.element.value).toBe(search);
  });

  it('Search hints was applied', () => {
    const { tooltipText } = mockData;
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
    const { tooltipText } = mockData;
    const tooltip = `
      <span class="tooltip-content-span">${tooltipText}</span>
    `;
    const wrapper = factory({
      propsData: {
        query: {},
        tooltip,
      },
    });
    const tooltipSpanElement = wrapper.find('.tooltip-content-span');

    expect(tooltipSpanElement.html()).toContain(tooltipText);
  });

  it('Submit search button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { query: { search } } });
    const input = wrapper.find('input');

    input.trigger('keydown');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
  });

  it('Submit search with search by column without columns', () => {
    const { column, search } = mockData;
    const advancedSearch = `- ${column.text} = '${search}'`;

    const wrapper = factory({
      propsData: {
        query: { search: advancedSearch },
      },
    });
    const submitButton = wrapper.find('.c-search-field_submit');

    submitButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
      search: `${column.text} = '${search}'`,
    }]);
  });

  it('Submit search with search by column', () => {
    const { column, search } = mockData;
    const advancedSearch = `- ${column.text} = '${search}'`;

    const wrapper = factory({
      propsData: {
        query: { search: advancedSearch },
        columns: [column],
      },
    });
    const submitButton = wrapper.find('.c-search-field_submit');

    submitButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
      search: `${column.value} = '${search}'`,
    }]);
  });

  it('Clear search by click on button', () => {
    const { search } = mockData;
    const wrapper = factory({
      propsData: {
        query: { search },
      },
    });
    const clearButton = wrapper.find('.c-search-field_clear');

    clearButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
    }]);
  });

  it('Submit search with custom field name', () => {
    const { search, field } = mockData;
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
    const { search, field } = mockData;
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
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
    }]);
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

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
