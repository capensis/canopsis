import Vue from 'vue';
import Vuetify from 'vuetify';
import { shallowMount } from '@vue/test-utils';

import { EXPAND_DEFAULT_MAX_LETTERS } from '@/config';

import Ellipsis from '@/components/tables/ellipsis.vue';

Vue.use(Vuetify);

describe('Ellipsis.vue', () => {
  it('displays default message', () => {
    const text = 'test';

    const wrapper = shallowMount(Ellipsis, {
      propsData: { text },
    });

    expect(wrapper.text()).toBe(text);
  });

  it('long text', () => {
    const text = '123456789012345678901234567890123456789012345678901234567890'; // 70 symbols
    const shortenText = `${text.substr(0, EXPAND_DEFAULT_MAX_LETTERS)}...${text}`;

    const wrapper = shallowMount(Ellipsis, {
      propsData: { text },
    });

    expect(wrapper.text()).toBe(shortenText);
  });
});
