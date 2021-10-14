import Faker from 'faker';

import { mount, createVueInstance } from '@unit/utils/vue';

import CSearchField from '@/components/forms/fields/c-search-field.vue';

const localVue = createVueInstance();

const mockData = {
  search: Faker.lorem.words(),
  newSearch: Faker.lorem.words(),
};

const stubs = {
  'v-btn': {
    template: `
      <button
        class="v-btn"
        @click="$listeners.click"
      >
        <slot />
      </button>
    `,
  },
  'v-text-field': {
    props: ['value'],
    template: `
      <input
        class="v-text-field"
        :value="value"
        @input="$emit('input', $event.target.value)"
        @keydown="$listeners.keydown"
      />
    `,
  },
};

const factory = (options = {}) => mount(CSearchField, {
  localVue,
  stubs,
  ...options,
});

describe('c-search-field', () => {
  it('Not empty value was pass into props and it was applied to input field', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = wrapper.find('input.v-text-field');

    expect(input.element.value).toBe(search);
  });

  it('Input value was update, after prop change and it was applied to input field', async () => {
    const { search, newSearch } = mockData;

    const wrapper = factory({ propsData: { value: search } });

    wrapper.setProps({ value: newSearch });

    await localVue.nextTick();

    const input = wrapper.find('input.v-text-field');

    expect(input.element.value).toBe(newSearch);
  });

  it('Set value into input element', () => {
    const { search } = mockData;

    const wrapper = factory();
    const input = wrapper.find('input.v-text-field');

    input.setValue(search);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([search]);
  });

  it('Keyup without enter key on input element', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = wrapper.find('input.v-text-field');

    input.trigger('keyup');

    const submitEvents = wrapper.emitted('submit');

    expect(submitEvents).toBeUndefined();
  });

  it('Keyup with enter key on input element', async () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = wrapper.find('input.v-text-field');

    await input.trigger('keydown.enter');

    const submitEvents = wrapper.emitted('submit');

    expect(submitEvents).toHaveLength(1);
  });

  it('Submit search button is the first button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const submitButton = wrapper.findAll('.v-btn').at(0);
    const submitIcon = submitButton.find('v-icon-stub');

    expect(submitIcon).toBeTruthy();
    expect(submitIcon.text()).toBe('search');
  });

  it('Clear search button is the second button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const clearButton = wrapper.findAll('.v-btn').at(1);

    const clearIcon = clearButton.find('v-icon-stub');

    expect(clearIcon).toBeTruthy();
    expect(clearIcon.text()).toBe('clear');
  });

  it('Click on submit button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const submitButton = wrapper.findAll('.v-btn').at(0);

    submitButton.trigger('click');

    const submitEvents = wrapper.emitted('submit');

    expect(submitEvents).toHaveLength(1);
  });

  it('Click on clear button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const clearButton = wrapper.findAll('.v-btn').at(1);

    clearButton.trigger('click');

    const inputEvents = wrapper.emitted('input');
    const clearEvents = wrapper.emitted('clear');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual(['']);

    expect(clearEvents).toHaveLength(1);
  });

  it('Renders `c-search-field` correctly', () => {
    const wrapper = mount(CSearchField, {
      localVue,
      propsData: { value: 'c-search-field' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
