import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import CEllipsis from '@/components/common/table/c-ellipsis.vue';

const stubs = {
  'v-menu': createActivatorElementStub('v-menu'),
  'c-simple-tooltip': true,
};

const hoverTooltipStub = {
  name: 'c-simple-tooltip',
  props: ['content', 'disabled', 'top'],
  template: `
    <div class="c-simple-tooltip">
      <slot name="activator" :on="{}" />
    </div>
  `,
};

const hoverStubs = {
  'v-menu': createActivatorElementStub('v-menu'),
  'c-simple-tooltip': hoverTooltipStub,
};

const defaultMaxLetters = CEllipsis.props.maxLetters.default;

const selectMenu = wrapper => wrapper.find('.v-menu');
const selectSpan = wrapper => selectMenu(wrapper).find('span');
const selectHoverTooltip = wrapper => wrapper.findComponent({ name: 'c-simple-tooltip' });
const selectHoverActivator = wrapper => wrapper.find('.c-ellipsis__activator');

describe('c-ellipsis', () => {
  const factory = generateShallowRenderer(CEllipsis, {
    stubs,
    attachTo: document.body,
  });

  it('Text letters count less then default max letters', () => {
    const text = Faker.datatype.string(defaultMaxLetters - 1);

    const wrapper = factory({ propsData: { text } });

    expect(wrapper.text()).toBe(text);
    expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
  });

  it('Text letters count more then default max letters', async () => {
    const text = Faker.datatype.string(defaultMaxLetters + 1);
    const shortenText = text.substr(0, defaultMaxLetters);

    const wrapper = factory({ propsData: { text } });

    await flushPromises();

    expect(wrapper.find('div > span').text()).toBe(shortenText);
    expect(selectSpan(wrapper).text()).toBe('...');
    expect(selectMenu(wrapper).exists()).toBeTruthy();
  });

  it('Text letters count less then custom maxLetters', () => {
    const text = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        text,
        maxLetters: text.length + 1,
      },
    });

    expect(wrapper.text()).toBe(text);
    expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
  });

  it('Text letters count more then custom maxLetters', () => {
    const maxLetters = Faker.datatype.number({ min: 1, max: 50 });
    const text = Faker.datatype.string(maxLetters + 1);
    const shortenText = text.substr(0, maxLetters);

    const wrapper = factory({ propsData: { text, maxLetters } });

    expect(wrapper.find('div > span').text()).toBe(shortenText);
    expect(selectSpan(wrapper).text()).toBe('...');
    expect(selectMenu(wrapper).exists()).toBeTruthy();
  });

  it('Click on text with text letters count less then default max letters', () => {
    const text = Faker.datatype.string();

    const wrapper = factory({ propsData: { text } });

    wrapper.find('div > span').trigger('click');

    expect(wrapper).toHaveBeenEmit('textClicked');
  });

  it('Click on text with text letters count more then default max letters', () => {
    const text = Faker.datatype.string();

    const wrapper = factory({ propsData: { text } });

    wrapper.find('div > span').trigger('click');

    expect(wrapper).toHaveBeenEmit('textClicked');
  });

  describe('hover', () => {
    const hoverFactory = generateShallowRenderer(CEllipsis, {
      stubs: hoverStubs,
      attachTo: document.body,
    });

    it('Text letters count less then default max letters', () => {
      const text = Faker.datatype.string(defaultMaxLetters - 1);

      const wrapper = hoverFactory({ propsData: { text, hover: true } });

      expect(selectHoverActivator(wrapper).text()).toBe(text);
      expect(selectHoverTooltip(wrapper).props()).toMatchObject({
        content: text,
        disabled: true,
      });
      expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
    });

    it('Text letters count more then default max letters', async () => {
      const text = Faker.datatype.string(defaultMaxLetters + 1);
      const shortenText = text.substr(0, defaultMaxLetters);

      const wrapper = hoverFactory({ propsData: { text, hover: true } });

      await flushPromises();

      expect(selectHoverActivator(wrapper).text().replace(/\s/g, '')).toBe(`${shortenText}...`);
      expect(selectHoverTooltip(wrapper).props()).toMatchObject({
        content: text,
        disabled: false,
      });
      expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
    });

    it('Text letters count less then custom maxLetters', () => {
      const text = Faker.datatype.string();

      const wrapper = hoverFactory({
        propsData: {
          text,
          hover: true,
          maxLetters: text.length + 1,
        },
      });

      expect(selectHoverActivator(wrapper).text()).toBe(text);
      expect(selectHoverTooltip(wrapper).props('disabled')).toBe(true);
      expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
    });

    it('Text letters count more then custom maxLetters', () => {
      const maxLetters = Faker.datatype.number({ min: 1, max: 50 });
      const text = Faker.datatype.string(maxLetters + 1);
      const shortenText = text.substr(0, maxLetters);

      const wrapper = hoverFactory({ propsData: { text, hover: true, maxLetters } });

      expect(selectHoverActivator(wrapper).text().replace(/\s/g, '')).toBe(`${shortenText}...`);
      expect(selectHoverTooltip(wrapper).props()).toMatchObject({
        content: text,
        disabled: false,
      });
    });

    it('Click on activator emits textClicked', () => {
      const text = Faker.datatype.string();

      const wrapper = hoverFactory({ propsData: { text, hover: true } });

      selectHoverActivator(wrapper).trigger('click');

      expect(wrapper).toHaveBeenEmit('textClicked');
    });
  });
});
