import { flushPromises, mount } from '@vue/test-utils';
import { createPinia } from 'pinia';
import { h } from 'vue';

import { Handlebars, unregisterHelper } from '@/helpers/handlebars';
import { sanitizeVueTemplate } from '@/helpers/html';

import { handlebarsLinksHelperCreator } from '@/mixins/handlebars/links-helper-creator';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';

const mixin = handlebarsLinksHelperCreator('entity.links', 'permission');
const LinksList = {
  props: { links: Object, category: { type: String, default: null } },
  render: () => h('div'),
};

const renderHelper = async (template, context = {}) => {
  const render = Handlebars.compile(template);
  const html = sanitizeVueTemplate(await render(context));
  const wrapper = mount(CRuntimeTemplate, {
    props: {
      template: html,
      templateProps: { links: {}, hasAccessForLinks: true },
    },
    global: { components: { 'c-links-list': LinksList } },
  });

  return wrapper;
};

describe('links Handlebars helper category', () => {
  beforeEach(() => {
    mixin.beforeCreate();
  });

  afterEach(() => {
    unregisterHelper('links');
    delete window.categoryXss;
  });

  it.each([
    ['{{links category="test"}}', {}, 'test'],
    ['{{links category=entity.category}}', { entity: { category: 'dynamic' } }, 'dynamic'],
    ['{{links category=entity.category}}', { entity: { category: 'nameOfVueVariable' } }, 'nameOfVueVariable'],
    ['{{links category=42}}', {}, '42'],
  ])('passes the resolved value as a string for %s', async (template, context, category) => {
    const wrapper = await renderHelper(template, context);

    expect(wrapper.findComponent(LinksList).props('category')).toBe(category);
    wrapper.unmount();
  });

  it.each([undefined, null, '', false, 0])('keeps the component default when category is %s', async (category) => {
    const wrapper = await renderHelper('{{links category=value}}', { value: category });

    expect(wrapper.findComponent(LinksList).props('category')).toBeNull();
    wrapper.unmount();
  });

  it('keeps the component default when the category argument is omitted', async () => {
    const wrapper = await renderHelper('{{links}}');

    expect(wrapper.findComponent(LinksList).props('category')).toBeNull();
    wrapper.unmount();
  });

  it('updates a category supplied by a Handlebars variable when the context changes', async () => {
    const wrapper = mount(CCompiledTemplate, {
      props: {
        template: '{{links category=entity.category}}',
        context: { entity: { category: 'first' } },
        templateProps: { links: {}, hasAccessForLinks: true },
      },
      global: {
        plugins: [createPinia()],
        components: { CRuntimeTemplate, 'c-links-list': LinksList },
        provide: { $system: { theme: {}, dark: false } },
      },
    });

    await flushPromises();
    expect(wrapper.findComponent(LinksList).props('category')).toBe('first');

    await wrapper.setProps({ context: { entity: { category: 'second' } } });
    await flushPromises();
    expect(wrapper.findComponent(LinksList).props('category')).toBe('second');
    wrapper.unmount();
  });

  it.each([
    'O\'Reilly "quoted" & <category>',
    '\' + [].filter.constructor("window.categoryXss = true")() + \'',
    new Handlebars.SafeString('" v-if="[].filter.constructor(\'window.categoryXss = true\')()'),
  ])('passes quoted and expression-like values as literal text', async (category) => {
    const wrapper = await renderHelper('{{links category=value}}', { value: category });

    expect(wrapper.findComponent(LinksList).props('category')).toBe(String(category));
    expect(window.categoryXss).toBeUndefined();
    wrapper.unmount();
  });
});
