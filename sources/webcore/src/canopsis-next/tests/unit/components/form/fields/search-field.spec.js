import Vuetify from 'vuetify';
import { generate as generateString } from 'randomstring';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import SearchField from '@/components/forms/fields/search-field.vue';

const localVue = createLocalVue();

const mocks = {
  $t: () => {},
};

const stubs = {
  'v-btn': {
    template: '<button class="v-btn" @click="$listeners.click"><slot></slot></button>',
  },
  'v-text-field': {
    props: ['value'],
    template: '<input class="v-text-field" :value="value" @input="$listeners.input($event.target.value)" @keyup="$listeners.keyup" />',
  },
};

localVue.use(Vuetify);

const factory = (options = {}) => shallowMount(SearchField, {
  localVue, mocks, stubs, ...options,
});

describe('SearchField', () => {
  it('Not empty value was pass into props and it was applied to input field', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const input = wrapper.find('input.v-text-field');

    expect(input.element.value).toBe(value);
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

  it('Keyup with enter key on input element', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const input = wrapper.find('input.v-text-field');

    input.trigger('keyup.enter');

    expect(wrapper.emitted('submit')).toBeTruthy();
    expect(wrapper.emitted('submit').length).toBe(1);
  });

  it('Submit search button is the second element', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const submitButton = wrapper.find('.v-btn:nth-child(2)');
    const submitIcon = submitButton.find('v-icon-stub');

    expect(submitIcon).toBeTruthy();
    expect(submitIcon.text()).toBe('search');
  });

  it('Clear search button is the third element', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const clearButton = wrapper.find('.v-btn:nth-child(3)');
    const clearIcon = clearButton.find('v-icon-stub');

    expect(clearIcon).toBeTruthy();
    expect(clearIcon.text()).toBe('clear');
  });

  it('Click on submit button', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const submitButton = wrapper.find('.v-btn:nth-child(2)');

    submitButton.trigger('click');

    expect(wrapper.emitted('submit')).toBeTruthy();
    expect(wrapper.emitted('submit').length).toBe(1);
  });

  it('Click on clear button', () => {
    const value = generateString();

    const wrapper = factory({ propsData: { value } });
    const clearButton = wrapper.find('.v-btn:nth-child(3)');

    clearButton.trigger('click');

    expect(wrapper.emitted('input')).toBeTruthy();
    expect(wrapper.emitted('input').length).toBe(1);
    expect(wrapper.emitted('input')[0]).toEqual(['']);

    expect(wrapper.emitted('clear')).toBeTruthy();
    expect(wrapper.emitted('clear').length).toBe(1);
    expect(wrapper.emittedByOrder().map(e => e.name)).toEqual(['input', 'clear']);
  });
});
