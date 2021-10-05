import Vuetify from 'vuetify';
import { generate as generateString } from 'randomstring';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import CAdvancedSearchField from '@/components/forms/fields/c-advanced-search-field.vue';

const localVue = createLocalVue();

const mocks = {
  $t: () => {},
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
  'v-tooltip': {
    template: `
      <div class='v-tooltip'>
        <slot />
      </div>
    `,
  },
};

localVue.use(Vuetify);

const factory = (options = {}) => shallowMount(CAdvancedSearchField, {
  localVue,
  mocks,
  stubs,
  ...options,
});

describe('c-advanced-search-field', () => {
  it('Pass default search into props and it was applied to input field', () => {
    const search = generateString();

    const wrapper = factory({ propsData: { query: { search } } });
    const input = wrapper.find('input');

    expect(input.element.value).toBe(search);
  });

  it('Search hints was applied', () => {
    const tooltipText = generateString();
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
    const tooltipContentText = generateString();
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
    const search = generateString();

    const wrapper = factory({ propsData: { query: { search } } });
    const input = wrapper.find('input');

    input.trigger('keydown');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
  });

  it('Submit search with search by column without columns', () => {
    const nameColumn = {
      text: 'Name',
      value: 'name',
    };
    const searchNameValue = generateString();
    const search = `- ${nameColumn.text} = '${searchNameValue}'`;

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
      search: `${nameColumn.text} = '${searchNameValue}'`,
    }]);
  });

  it('Submit search with search by column', () => {
    const nameColumn = {
      text: 'Name',
      value: 'name',
    };
    const searchNameValue = generateString();
    const search = `- ${nameColumn.text} = '${searchNameValue}'`;

    const wrapper = factory({
      propsData: {
        query: { search },
        columns: [nameColumn],
      },
    });
    const submitButton = wrapper.find('.c-search-field_submit');

    submitButton.trigger('click');

    const updateQueryEvents = wrapper.emitted('update:query');

    expect(updateQueryEvents).toHaveLength(1);
    expect(updateQueryEvents[0]).toEqual([{
      page: 1,
      search: `${nameColumn.value} = '${searchNameValue}'`,
    }]);
  });

  it('Clear search by click on button', () => {
    const search = generateString();
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
    const search = generateString();
    const field = generateString();
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
    const search = generateString();
    const field = generateString();
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
    const wrapper = shallowMount(CAdvancedSearchField, {
      localVue,
      mocks,
      stubs: ['c-search-field'],
      propsData: {
        query: { search: 'c-advanced-search-field' },
        tooltip,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
