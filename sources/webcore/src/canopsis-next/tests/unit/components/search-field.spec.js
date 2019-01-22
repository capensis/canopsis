import Vuetify from 'vuetify';
import { generate as generateString } from 'randomstring';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import SearchField from '@/components/forms/fields/search-field.vue';

const localVue = createLocalVue();

const mocks = {
  $t: () => {},
};

const stubs = {
  'v-text-field': {
    props: ['value'],
    template: '<input class="v-text-field" :value="value" @input="$listeners.input($event.target.value)" @keyup.enter="$listeners.keyup" />',
  },
};

localVue.use(Vuetify);

describe('SearchField', () => {
  it('Text letters count less then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const value = generateString();
    const wrapper = shallowMount(SearchField, {
      localVue,
      mocks,
      stubs,
    });

    const input = wrapper.find('.v-text-field');

    input.setValue(value);

    expect(wrapper.emitted('input')).toBeTruthy();
    expect(wrapper.emitted('input').length).toBe(1);
    expect(wrapper.emitted('input')[0]).toEqual([value]);
  });

  it('Text letters count less then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const value = generateString();
    const wrapper = shallowMount(SearchField, {
      localVue,
      mocks,
      stubs,
    });

    const input = wrapper.find('.v-text-field');

    input.setValue(value);
    input.trigger('keyup');

    expect(wrapper.emitted('submit')).toBeUndefined();
  });

  it('Text letters count less then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const value = generateString();
    const wrapper = shallowMount(SearchField, {
      localVue,
      mocks,
      stubs,
    });

    const input = wrapper.find('.v-text-field');

    input.setValue(value);
    input.trigger('keyup.enter');

    expect(wrapper.emitted('submit')).toBeTruthy();
    expect(wrapper.emitted('submit').length).toBe(1);
  });
});
