import { flushPromises, mount } from '@vue/test-utils';
import { createPinia } from 'pinia';
import { h } from 'vue';

import { vuetify } from '@/plugins/vuetify';

import { Handlebars } from '@/helpers/handlebars';
import { alarmOrEntityTagsHelper } from '@/helpers/handlebars/tags-helper';
import { sanitizeVueTemplate } from '@/helpers/html';

import CIcon from '@/components/common/icon/c-icon.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';

const payload = 'constructor.constructor(\'window.templateXss = true\')()';
const componentStub = {
  props: ['alarm', 'selectedTags', 'links', 'category', 'service', 'options'],
  emits: ['select', 'close', 'refresh', 'update:options', 'add:action'],
  render() {
    return h('span', this.$slots.default?.());
  },
};

const createParentContext = (values = {}) => ({
  ...values,
  $: { provides: {} },
  $data: Object.fromEntries(Object.entries(values).filter(([, value]) => typeof value !== 'function')),
  $options: {
    methods: Object.fromEntries(Object.entries(values).filter(([, value]) => typeof value === 'function')),
  },
});

const factory = (template, parent = {}, { sanitizeOptions, ...options } = {}) => mount(CRuntimeTemplate, {
  props: {
    ...options,
    template: sanitizeVueTemplate(template, sanitizeOptions),
    parent: createParentContext(parent),
  },
  global: {
    plugins: [vuetify],
    components: {
      CIcon,
      'c-alarm-tags-chips': componentStub,
      'c-links-list': componentStub,
      'service-entities-list': componentStub,
    },
  },
});

describe('sanitized HTML in the unchanged c-runtime-template', () => {
  afterEach(() => {
    delete window.templateXss;
  });

  it('retains the original runtime compiler behavior for trusted templates', () => {
    const wrapper = mount(CRuntimeTemplate, {
      props: { template: '<span>{{ value }}</span>', templateProps: { value: 'trusted content' } },
    });

    expect(wrapper.text()).toBe('trusted content');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('renders safe comments as DOM comments', () => {
    const wrapper = factory('<div><!-- Комментарий шаблона -->text</div>');

    expect(wrapper.html()).toContain('<!-- Комментарий шаблона -->');
    expect(wrapper.text()).toBe('text');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it.each([
    `<span v-if="$root.${payload}">text</span>`,
    `<span v-else-if="$root.${payload}">text</span>`,
    `<span v-bind:title="${payload}">text</span>`,
    `<span :title="${payload}">text</span>`,
    `<span @click="${payload}">text</span>`,
    `<span v-on:click="${payload}">text</span>`,
    `<span>{{ ${payload} }}</span>`,
    `<span v-html="${payload}">text</span>`,
    `<span data-v-on-click="${payload}">text</span>`,
    `<span data-v-bind-title="${payload}">text</span>`,
    `<span v-text="${payload}">text</span>`,
    `<span v-for="item in ${payload}">text</span>`,
    `<span v-on="${payload}">text</span>`,
    `<span v-bind="${payload}">text</span>`,
    `<span :[${payload}]="1">text</span>`,
    `<span v-slot="${payload}">text</span>`,
  ])('does not interpret Vue syntax: %s', async (template) => {
    const wrapper = factory(template);

    await wrapper.find('span').trigger('click');

    expect(wrapper.find('span').attributes()).toEqual({});
    expect(wrapper.text()).toBe(template.includes('{{') ? `{{ ${payload} }}` : 'text');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('sanitizes native handlers, URLs and component DOM override props', () => {
    const wrapper = factory(`<div><img src="x" onerror="${payload}">
      <a href="javascript:${payload}">link</a>
      <v-icon tag="script" :innerHTML="${payload}" onload="${payload}">home</v-icon>
      </div>`);

    expect(wrapper.find('img').attributes()).toEqual({ src: 'x' });
    expect(wrapper.find('a').attributes('href')).toBeUndefined();
    expect(wrapper.find('script').exists()).toBe(false);
    expect(wrapper.find('i').text()).toBe('home');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it.each(['v-icon', 'c-icon'].flatMap(tag => [
    `<${tag} @click="${payload}">home</${tag}>`,
    `<${tag} v-on:click="${payload}">home</${tag}>`,
    `<${tag} data-v-on-click="${payload}">home</${tag}>`,
    `<${tag} :size="${payload}">home</${tag}>`,
    `<${tag} v-bind="${payload}">home</${tag}>`,
    `<${tag} v-html="${payload}">home</${tag}>`,
    `<${tag} tag="script" is="script" innerHTML="<img src=x onerror=alert(1)>">home</${tag}>`,
    `<${tag}>{{ ${payload} }}</${tag}>`,
  ]))('does not execute untrusted icon attributes or interpolation: %s', async (template) => {
    const wrapper = factory(template);
    const icon = wrapper.find('i');

    expect(wrapper.findComponent(CIcon).exists()).toBe(true);
    await icon.trigger('click');
    await icon.trigger('load');
    expect(wrapper.find('script, img').exists()).toBe(false);
    expect(icon.attributes('tag')).toBeUndefined();
    expect([undefined, '']).toContain(icon.attributes('is'));
    expect(wrapper.text()).toBe(template.includes('{{') ? `{{ ${payload} }}` : 'home');
    expect(window.templateXss).toBeUndefined();
    wrapper.unmount();
  });

  it('does not instantiate arbitrary registered components', () => {
    const mounted = vi.fn();
    const wrapper = mount(CRuntimeTemplate, {
      props: { template: sanitizeVueTemplate('<dangerous-component>text</dangerous-component>') },
      global: { components: { 'dangerous-component': { mounted, template: '<div />' } } },
    });

    expect(mounted).not.toHaveBeenCalled();
    expect(wrapper.text()).toBe('text');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('preserves rich HTML, link attributes, SVG and literal Vue interpolation', () => {
    const wrapper = factory('<section><b>bold</b><a href="https://example.com">link</a><svg viewBox="0 0 24 24"><path d="M0 0" /></svg>{{ literal }}</section>');

    expect(wrapper.find('section').exists()).toBe(true);
    expect(wrapper.find('b').text()).toBe('bold');
    expect(wrapper.find('a').attributes('href')).toBe('https://example.com');
    expect(wrapper.find('svg').element.namespaceURI).toBe('http://www.w3.org/2000/svg');
    expect(wrapper.text()).toContain('{{ literal }}');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('preserves custom sanitization options such as text widget iframes', () => {
    const wrapper = factory('<div><iframe src="https://example.com"></iframe></div>', {}, {
      sanitizeOptions: { ADD_TAGS: ['iframe'] },
    });

    expect(wrapper.find('iframe').attributes('src')).toBe('https://example.com');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('supports the fixed tag helper bindings and event forwarding', async () => {
    const alarm = { _id: 'alarm-id' };
    const selectedTags = ['tag'];
    const select = vi.fn();
    const close = vi.fn();
    const template = String(alarmOrEntityTagsHelper.call({ alarm }, { hash: {} }));
    const wrapper = factory(template, { alarm, selectedTags }, {
      'onSelect:tag': select,
      'onRemove:tag': close,
    });
    const tags = wrapper.findComponent(componentStub);

    expect(tags.props('alarm')).toEqual(alarm);
    expect(tags.props('selectedTags')).toEqual(selectedTags);
    tags.vm.$emit('select', 'tag');
    tags.vm.$emit('close', 'tag');
    await flushPromises();
    expect(select).toHaveBeenCalledWith('tag');
    expect(close).toHaveBeenCalledWith('tag');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('does not evaluate expressions in component bindings', () => {
    const wrapper = factory(`<c-alarm-tags-chips :alarm="${payload}" />`);

    expect(wrapper.findComponent(componentStub).props('alarm')).toBeUndefined();
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it.each([
    `{{ ${payload} }}`,
    `&#123;&#123; ${payload} &#125;&#125;`,
    `&#x7b;&#x7b; ${payload} }}`,
    `&lbrace;{ ${payload} }}`,
    `{{{ ${payload} }}}`,
    '{{ [].filter.constructor(\'window.templateXss = true\')() }}',
  ])('keeps interpolation text inert after HTML entity normalization: %s', (text) => {
    const wrapper = factory(`<span>${text}</span>`);

    expect(wrapper.text()).toContain('window.templateXss = true');
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('removes event handlers whose values are not the fixed helper callbacks', () => {
    const unexpected = vi.fn();
    const wrapper = factory('<service-entities-list @refresh="unexpected"></service-entities-list>', { unexpected });

    wrapper.findComponent(componentStub).vm.$emit('refresh', true);
    expect(unexpected).not.toHaveBeenCalled();
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('removes bindings that execute code without using the parent context', () => {
    const wrapper = factory('<span :title="[].filter.constructor(\'window.templateXss = true\')()">text</span>');

    expect(wrapper.find('span').attributes()).toEqual({});
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('rejects Vue property shorthand even when custom HTML options allow the attribute', () => {
    const wrapper = factory('<span .title="[].filter.constructor(\'window.templateXss = true\')()">text</span>', {}, {
      sanitizeOptions: { ADD_ATTR: ['.title'] },
    });

    expect(wrapper.find('span').attributes()).toEqual({});
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('enforces the links helper access check and preserves literal categories', async () => {
    const parent = { links: [{ label: 'link' }], hasAccessForLinks: false };
    const template = '<c-links-list v-if="links &amp;&amp; hasAccessForLinks" :links="links" category="test"></c-links-list>';
    const wrapper = factory(template, parent);

    expect(wrapper.findComponent(componentStub).exists()).toBe(false);
    await wrapper.setProps({ parent: createParentContext({ ...parent, hasAccessForLinks: true }) });
    expect(wrapper.findComponent(componentStub).props()).toMatchObject({
      links: parent.links,
      category: 'test',
    });
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('forwards only the fixed service helper events', () => {
    const refreshEntities = vi.fn();
    const unexpected = vi.fn();
    const wrapper = factory('<service-entities-list @refresh="refreshEntities"></service-entities-list>', {
      refreshEntities,
      unexpected,
    });

    wrapper.findComponent(componentStub).vm.$emit('refresh', true);
    expect(refreshEntities).toHaveBeenCalledWith(true);
    expect(unexpected).not.toHaveBeenCalled();
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('keeps untrusted content inert through Handlebars and the compiled renderer', async () => {
    const wrapper = mount(CCompiledTemplate, {
      props: {
        template: '{{{ comment }}}',
        context: { comment: `<span v-if="$root.${payload}">comment</span>` },
      },
      global: {
        plugins: [createPinia()],
        components: { 'c-runtime-template': CRuntimeTemplate },
        provide: { $system: { theme: {}, dark: false } },
      },
    });

    await flushPromises();
    expect(wrapper.text()).toBe('comment');
    expect(wrapper.find('span').attributes()).toEqual({});
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
  });

  it('does not interpret Vue interpolation introduced by Handlebars data', async () => {
    Handlebars.registerHelper('vueLiteral', () => new Handlebars.SafeString(`{{ ${payload} }}`));
    const wrapper = mount(CCompiledTemplate, {
      props: { template: '{{vueLiteral}}' },
      global: {
        plugins: [createPinia()],
        components: { 'c-runtime-template': CRuntimeTemplate },
        provide: { $system: { theme: {}, dark: false } },
      },
    });

    await flushPromises();
    expect(wrapper.text()).toBe(`{{ ${payload} }}`);
    wrapper.unmount();
    expect(window.templateXss).toBeUndefined();
    Handlebars.unregisterHelper('vueLiteral');
  });
});
