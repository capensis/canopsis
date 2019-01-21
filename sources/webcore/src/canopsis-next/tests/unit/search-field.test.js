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
    template: '<input class="v-text-field" :value="value" @input="$listeners.input" @keyup.enter="$listeners.keyup" />',
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

    expect(wrapper.emitted('input')).toBeDefined();
    expect(wrapper.emitted('input').length).toBe(1);

    input.trigger('keyup.enter');

    expect(wrapper.emitted('submit')).toBeDefined();
  });
});
