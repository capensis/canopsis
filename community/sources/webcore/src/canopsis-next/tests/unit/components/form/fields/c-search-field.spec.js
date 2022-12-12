import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CSearchField from '@/components/forms/fields/c-search-field.vue';

const localVue = createVueInstance();

const mockData = {
  search: Faker.lorem.words(),
  newSearch: Faker.lorem.words(),
};

const stubs = {
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
  'c-action-btn': true,
};

const snapshotStubs = {
  'c-action-btn': true,
};

describe('c-search-field', () => {
  const factory = generateShallowRenderer(CSearchField, {
    localVue,
    stubs,
  });
  const snapshotFactory = generateRenderer(CSearchField, {
    localVue,
    stubs: snapshotStubs,
  });

  it('Not empty value was pass into props and it was applied to input field', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = wrapper.find('input.v-text-field');

    expect(input.element.value).toBe(search);
  });

  it('Input value was update, after prop change and it was applied to input field', async () => {
    const { search, newSearch } = mockData;

    const wrapper = factory({ propsData: { value: search } });

    await wrapper.setProps({ value: newSearch });

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
    const submitButton = wrapper.findAll('c-action-btn-stub').at(0);

    expect(submitButton.attributes('icon')).toBe('search');
  });

  it('Clear search button is the second button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const clearButton = wrapper.findAll('c-action-btn-stub').at(1);

    expect(clearButton.attributes('icon')).toBe('clear');
  });

  it('Click on submit button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const submitButton = wrapper.findAll('c-action-btn-stub').at(0);

    submitButton.vm.$emit('click');

    const submitEvents = wrapper.emitted('submit');

    expect(submitEvents).toHaveLength(1);
  });

  it('Click on clear button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const clearButton = wrapper.findAll('c-action-btn-stub').at(1);

    clearButton.vm.$emit('click');

    const inputEvents = wrapper.emitted('input');
    const clearEvents = wrapper.emitted('clear');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual(['']);

    expect(clearEvents).toHaveLength(1);
  });

  it('Renders `c-search-field` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: { value: 'c-search-field' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
