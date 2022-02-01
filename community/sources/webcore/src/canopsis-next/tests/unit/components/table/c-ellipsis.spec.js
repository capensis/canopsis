import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CEllipsis from '@/components/common/table/c-ellipsis.vue';

const localVue = createVueInstance();

const defaultMaxLetters = CEllipsis.props.maxLetters.default;

const mockData = {
  shortText: Faker.datatype.string(defaultMaxLetters - 1),
  longText: Faker.datatype.string(defaultMaxLetters + 1),
  maxLetters: Faker.datatype.number(),
};

const factory = (options = {}) => shallowMount(CEllipsis, { localVue, ...options });

describe('c-ellipsis', () => {
  it('Text letters count less then default max letters', () => {
    const { shortText } = mockData;

    const wrapper = factory({ propsData: { text: shortText } });

    expect(wrapper.text()).toBe(shortText);
    expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
  });

  it('Text letters count more then default max letters', () => {
    const { longText } = mockData;
    const shortenText = longText.substr(0, defaultMaxLetters);

    const wrapper = factory({ propsData: { text: longText } });

    expect(wrapper.find('div > span').text()).toBe(shortenText);
    expect(wrapper.find('v-menu-stub > span').text()).toBe('...');
    expect(wrapper.find('v-menu-stub > v-card-stub > v-card-title-stub').text()).toBe(longText);
  });

  it('Text letters count less then custom maxLetters', () => {
    const { maxLetters, shortText } = mockData;

    const wrapper = factory({ propsData: { text: shortText, maxLetters } });

    expect(wrapper.text()).toBe(shortText);
    expect(wrapper.find('v-menu-stub').exists()).toBeFalsy();
  });

  it('Text letters count more then custom maxLetters', () => {
    const { maxLetters } = mockData;
    const text = Faker.datatype.string(maxLetters + 1);
    const shortenText = text.substr(0, maxLetters);

    const wrapper = factory({ propsData: { text, maxLetters } });

    expect(wrapper.find('div > span').text()).toBe(shortenText);
    expect(wrapper.find('v-menu-stub > span').text()).toBe('...');
    expect(wrapper.find('v-menu-stub > v-card-stub > v-card-title-stub').text()).toBe(text);
  });

  it('Click on dots with text letters count more then default max letters', async () => {
    const { longText } = mockData;

    const wrapper = factory({ propsData: { text: longText } });

    wrapper.find('v-menu-stub > span').trigger('click');

    await localVue.nextTick();

    expect(wrapper.vm.isFullTextMenuOpen).toBeTruthy();
  });

  it('Click on text with text letters count less then default max letters', () => {
    const { longText } = mockData;

    const wrapper = factory({ propsData: { text: longText } });

    wrapper.find('div > span').trigger('click');

    const textClickedEvents = wrapper.emitted('textClicked');

    expect(textClickedEvents).toHaveLength(1);
  });

  it('Click on text with text letters count more then default max letters', () => {
    const { longText } = mockData;

    const wrapper = factory({ propsData: { text: longText } });

    wrapper.find('div > span').trigger('click');

    const textClickedEvents = wrapper.emitted('textClicked');

    expect(textClickedEvents).toHaveLength(1);
  });

  it('Renders `c-ellipsis` correctly', () => {
    const text = `omnis esse recusandae magni similique porro quaerat alias vel deserunt porro voluptate
      voluptatibus commodi sequi et qui fugiat exercitationem et et eligendi ad officia quis suscipit earum soluta
      minima architecto numquam et voluptatibus quia officiis nulla veritatis soluta optio assumenda est fugiat est
      suscipit inventore temporibus dolores quos dolorem doloremque autem qui sit ipsam praesentium esse sunt ut
      molestiae nulla itaque voluptatem pariatur vel dolorum impedit asperiores numquam animi esse laborum in et
      magnam nulla et consequatur facere sint nam facere sunt aut alias qui omnis rerum corporis totam quibusdam
      nostrum mollitia quia vel amet pariatur eveniet explicabo quia ullam`;

    const wrapper = mount(CEllipsis, {
      localVue,
      propsData: { text },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
