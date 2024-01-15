import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createActivatorElementStub } from '@unit/stubs/vuetify';

import CEllipsis from '@/components/common/table/c-ellipsis.vue';

const stubs = {
  'v-menu': createActivatorElementStub('v-menu'),
};

const defaultMaxLetters = CEllipsis.props.maxLetters.default;

const selectMenu = wrapper => wrapper.find('.v-menu');
const selectSpan = wrapper => selectMenu(wrapper).find('span');
const selectCardTitle = wrapper => selectMenu(wrapper).find('v-card-stub > v-card-title-stub');

describe('c-ellipsis', () => {
  const factory = generateShallowRenderer(CEllipsis, {
    stubs,
    attachTo: document.body,
  });
  const snapshotFactory = generateRenderer(CEllipsis, {
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
    expect(selectCardTitle(wrapper).text()).toBe(text);
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
    expect(selectCardTitle(wrapper).text()).toBe(text);
  });

  it('Click on text with text letters count less then default max letters', () => {
    const text = Faker.datatype.string();

    const wrapper = factory({ propsData: { text } });

    wrapper.find('div > span').trigger('click');

    const textClickedEvents = wrapper.emitted('textClicked');

    expect(textClickedEvents).toHaveLength(1);
  });

  it('Click on text with text letters count more then default max letters', () => {
    const text = Faker.datatype.string();

    const wrapper = factory({ propsData: { text } });

    wrapper.find('div > span').trigger('click');

    const textClickedEvents = wrapper.emitted('textClicked');

    expect(textClickedEvents).toHaveLength(1);
  });

  it('Renders `c-ellipsis` correctly', async () => {
    const text = `omnis esse recusandae magni similique porro quaerat alias vel deserunt porro voluptate
      voluptatibus commodi sequi et qui fugiat exercitationem et et eligendi ad officia quis suscipit earum soluta
      minima architecto numquam et voluptatibus quia officiis nulla veritatis soluta optio assumenda est fugiat est
      suscipit inventore temporibus dolores quos dolorem doloremque autem qui sit ipsam praesentium esse sunt ut
      molestiae nulla itaque voluptatem pariatur vel dolorum impedit asperiores numquam animi esse laborum in et
      magnam nulla et consequatur facere sint nam facere sunt aut alias qui omnis rerum corporis totam quibusdam
      nostrum mollitia quia vel amet pariatur eveniet explicabo quia ullam`;

    const wrapper = snapshotFactory({
      propsData: { text },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
