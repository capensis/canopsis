import Vuetify from 'vuetify';
import { generate as generateString } from 'randomstring';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import CSearchField from '@/components/forms/fields/c-search-field.vue';

const localVue = createLocalVue();

const mocks = {
  $t: () => {},
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
        @keydown="$emit('keydown', $event)"
      />
    `,
  },
};

localVue.use(Vuetify);

const factory = (options = {}) => shallowMount(CSearchField, {
  localVue,
  mocks,
  stubs,
  ...options,
});

describe('c-search-field', () => {
  it('Not empty value was pass into props and it was applied to input field', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const input = wrapper.find('input.v-text-field');

    expect(input.element.value).toBe(value);
  });

  it('Input value was update, after prop change and it was applied to input field', async () => {
    const value = generateString();
    const newValue = generateString();

    const wrapper = factory({ propsData: { value } });

    wrapper.setProps({ value: newValue });

    await localVue.nextTick();

    const input = wrapper.find('input.v-text-field');

    expect(input.element.value).toBe(newValue);
  });

  it('Set value into input element', () => {
    const value = generateString();

    const wrapper = factory();
    const input = wrapper.find('input.v-text-field');

    input.setValue(value);

    expect(wrapper.emitted('input')).toBeTruthy();
    expect(wrapper.emitted('input').length).toBe(1);
    expect(wrapper.emitted('input')[0]).toEqual([value]);
  });

  it('Keyup without enter key on input element', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const input = wrapper.find('input.v-text-field');

    input.trigger('keyup');

    expect(wrapper.emitted('submit')).toBeUndefined();
  });

  it('Keyup with enter key on input element', async () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const input = wrapper.find('input.v-text-field');

    await input.trigger('keydown.enter');

    const submitEvents = wrapper.emitted('submit');

    expect(submitEvents).toBeTruthy();
    expect(submitEvents.length).toBe(1);
  });

  it('Submit search button is the first button', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const submitButton = wrapper.findAll('.v-btn').at(0);
    const submitIcon = submitButton.find('v-icon-stub');

    expect(submitIcon).toBeTruthy();
    expect(submitIcon.text()).toBe('search');
  });

  it('Clear search button is the second button', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const clearButton = wrapper.findAll('.v-btn').at(1);

    const clearIcon = clearButton.find('v-icon-stub');

    expect(clearIcon).toBeTruthy();
    expect(clearIcon.text()).toBe('clear');
  });

  it('Click on submit button', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const submitButton = wrapper.findAll('.v-btn').at(0);

    submitButton.trigger('click');

    expect(wrapper.emitted('submit')).toBeTruthy();
    expect(wrapper.emitted('submit').length).toBe(1);
  });

  it('Click on clear button', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const clearButton = wrapper.findAll('.v-btn').at(1);

    clearButton.trigger('click');

    const inputEvents = wrapper.emitted('input');
    const clearEvents = wrapper.emitted('clear');

    expect(inputEvents).toBeTruthy();
    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual(['']);

    expect(clearEvents).toBeTruthy();
    expect(clearEvents.length).toBe(1);
    expect(wrapper.emittedByOrder().map(e => e.name)).toEqual(['input', 'clear']);
  });

  it('Renders `c-search-field` correctly', () => {
    const wrapper = shallowMount(CSearchField, {
      localVue,
      mocks,
      propsData: { value: 'c-search-field' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
