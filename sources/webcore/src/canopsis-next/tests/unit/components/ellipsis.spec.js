import Vuetify from 'vuetify';
import { generate as generateString } from 'randomstring';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import { EXPAND_DEFAULT_MAX_LETTERS } from '../../../src/config';

import Ellipsis from '../../../src/components/tables/ellipsis.vue';

const localVue = createLocalVue();

localVue.use(Vuetify);

describe('Ellipsis', () => {
  it('Text letters count less then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const text = generateString(EXPAND_DEFAULT_MAX_LETTERS - 1);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text },
    });

    expect(wrapper.text()).toBe(text);
    expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
  });

  it('Text letters count more then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const text = generateString(EXPAND_DEFAULT_MAX_LETTERS + 1);
    const shortenText = text.substr(0, EXPAND_DEFAULT_MAX_LETTERS);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text },
    });

    expect(wrapper.find('div > span').text()).toBe(shortenText);
    expect(wrapper.find('v-menu-stub > span').text()).toBe('...');
    expect(wrapper.find('v-menu-stub > v-card-stub > v-card-title-stub').text()).toBe(text);
  });

  it('Text letters count less then custom maxLetters', () => {
    const maxLetters = 5;
    const text = generateString(maxLetters - 1);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text, maxLetters },
    });

    expect(wrapper.text()).toBe(text);
  });

  it('Text letters count more then custom maxLetters', () => {
    const maxLetters = 5;
    const text = generateString(maxLetters + 1);
    const shortenText = text.substr(0, maxLetters);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text, maxLetters },
    });

    expect(wrapper.find('div > span').text()).toBe(shortenText);
    expect(wrapper.find('v-menu-stub > span').text()).toBe('...');
    expect(wrapper.find('v-menu-stub > v-card-stub > v-card-title-stub').text()).toBe(text);
  });

  it('Click on dots with text letters count more then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const text = generateString(EXPAND_DEFAULT_MAX_LETTERS + 1);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text },
    });

    wrapper.find('v-menu-stub > span').trigger('click');

    expect(wrapper.vm.isFullTextMenuOpen).toBeTruthy();
  });

  it('Click on text with text letters count less then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const text = generateString(EXPAND_DEFAULT_MAX_LETTERS - 1);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text },
    });

    wrapper.find('div > span').trigger('click');

    expect(wrapper.emitted('textClicked')).toBeTruthy();
    expect(wrapper.emitted('textClicked').length).toBe(1);
  });

  it('Click on text with text letters count more then EXPAND_DEFAULT_MAX_LETTERS', () => {
    const text = generateString(EXPAND_DEFAULT_MAX_LETTERS + 1);

    const wrapper = shallowMount(Ellipsis, {
      localVue,
      propsData: { text },
    });

    wrapper.find('div > span').trigger('click');

    expect(wrapper.emitted('textClicked')).toBeTruthy();
    expect(wrapper.emitted('textClicked').length).toBe(1);
  });
});
