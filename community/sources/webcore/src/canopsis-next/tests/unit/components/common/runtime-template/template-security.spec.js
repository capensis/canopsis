import Vue from 'vue';
import VueRouter from 'vue-router';
import VueClipboard from 'vue-clipboard2';
import { mount } from '@vue/test-utils';
import axios from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';

import { flushPromises, generateRenderer, createVueInstance } from '@unit/utils/vue';

import i18n from '@/i18n';

import { sanitizeHtml, sanitizeVueTemplate } from '@/helpers/html';
import { registerTemplate, unregisterTemplate } from '@/helpers/handlebars';
import { alarmOrEntityTagsHelper } from '@/helpers/handlebars/tags-helper';

import { handlebarsLinksHelperCreator } from '@/mixins/handlebars/links-helper-creator';
import { vuetifyCustomIconsBaseMixin } from '@/mixins/vuetify/custom-icons/base';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CRequestHelper from '@/components/common/runtime-template/c-request-helper.vue';
import CAlarmTagsChips from '@/components/common/chips/c-alarm-tags-chips.vue';
import CEntityTagsChips from '@/components/common/chips/c-entity-tags-chips.vue';
import CAlarmActionsChips from '@/components/common/chips/c-alarm-actions-chips.vue';
import CChip from '@/components/common/chips/c-chip.vue';
import CPagination from '@/components/common/pagination/c-pagination.vue';
import CItemsPerPageField from '@/components/forms/fields/c-items-per-page-field.vue';
import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';
import CLinksList from '@/components/common/links/c-links-list.vue';
import CCopyWrapper from '@/components/common/runtime-template/c-copy-wrapper.vue';
import CNow from '@/components/common/runtime-template/c-now.vue';
import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';
import ServiceTemplate from '@/components/other/service/partials/service-template.vue';
import ServiceEntitiesList from '@/components/other/service/partials/service-entities-list.vue';
import CustomIcon from '@/components/other/icons/partials/custom-icon.vue';
import IconFileSelector from '@/components/other/icons/form/fields/icon-file-selector.vue';

const expression = "$root.constructor.constructor('window.templateXss = true')()";
const maliciousSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" onload="window.templateXss=true"><foreignObject width="1" height="1"><img xmlns="http://www.w3.org/1999/xhtml" src="x" onerror="window.templateXss=true"/></foreignObject><text>{{ ${expression} }}</text><path d="M0 0"/></svg>`;
const components = {
  CCompiledTemplate,
  CRuntimeTemplate,
  CRequestHelper,
  CAlarmChip,
  CChip,
  CPagination,
  CItemsPerPageField,
  CAlarmTagsChips,
  CEntityTagsChips,
  CAlarmActionsChips,
  CLinksList,
  CCopyWrapper,
  CNow,
  CTablePagination,
};
const localVue = createVueInstance();
localVue.use(VueClipboard);
localVue.use(VueRouter);
localVue.mixin({
  i18n: { locale: i18n.locale, fallbackLocale: i18n.fallbackLocale, messages: i18n.messages },
});
Object.entries(components).forEach(([name, component]) => localVue.component(name, component));
const parentComponent = { provide: { $system: { theme: {}, dark: false } } };

describe('actual Vue 2 runtime compilation after sanitization', () => {
  const factory = generateRenderer(CCompiledTemplate, { localVue, parentComponent });

  afterEach(() => {
    delete window.templateXss;
    unregisterTemplate('xss-regression');
  });

  it.each([false, true])('handles nested HBS conditions and iteration in table comments: %s', async (enabled) => {
    const source = '<!--{{#if enabled}}--><table><tbody><!--{{#each rows}}--><!--{{#if visible}}--><tr><td>{{name}}</td></tr><!--{{else}}--><tr><td>hidden</td></tr><!--{{/if}}--><!--{{/each}}--></tbody></table><!--{{else}}--><p>disabled</p><!--{{/if}}-->';
    const template = sanitizeHtml(sanitizeHtml(source));
    const wrapper = factory({ propsData: {
      template,
      context: { enabled, rows: [{ visible: true, name: 'yes' }, { visible: false }] },
    } });
    await flushPromises();

    expect(wrapper.text()).toBe(enabled ? 'yeshidden' : 'disabled');
    expect(wrapper.findAll('tbody > tr')).toHaveLength(enabled ? 2 : 0);
    await wrapper.setProps({ context: { enabled: true, rows: [{ visible: true, name: 'updated' }] } });
    await flushPromises();
    expect(wrapper.find('tbody > tr > td').text()).toBe('updated');
  });

  it.each([false, true])('rejects expressions and interpolations from HBS data, registered=%s', async (registered) => {
    const template = '{{{comment}}}';
    if (registered) {
      registerTemplate('xss-regression', template);
    }
    const wrapper = factory({ propsData: {
      template,
      templateId: registered ? 'xss-regression' : '',
      context: { comment: `<span v-if="${expression}">comment</span>&#123;&#123; ${expression} }}` },
    } });
    await flushPromises();

    expect(wrapper.text()).toBe(`comment{{ ${expression} }}`);
    expect(wrapper.find('span').attributes()).toEqual({});
    expect(window.templateXss).toBeUndefined();
  });

  it('checks parentElement together with the compiled content', async () => {
    const wrapper = factory({ propsData: { template: 'safe', parentElement: `div v-if="${expression}"` } });
    await flushPromises();

    expect(wrapper.text()).toContain('safe');
    expect(wrapper.find('[v-if]').exists()).toBe(false);
    expect(window.templateXss).toBeUndefined();
  });

  it('keeps internal links as RouterLink components', async () => {
    const router = new VueRouter({ mode: 'abstract', routes: [] });
    const wrapper = factory({ router, propsData: { template: '{{internal-link href="/alarms" text="Alarms"}}' } });
    await flushPromises();

    expect(wrapper.findComponent({ name: 'RouterLink' }).props('to')).toBe('/alarms');
    expect(wrapper.find('a').attributes('href')).toBe('/alarms');
  });

  it('keeps entity tag helper data', async () => {
    const entity = { tag_colors: [{ value: 'entity-tag', color: '#ff0000' }] };
    const template = String(alarmOrEntityTagsHelper.call({ entity }, { hash: {} }));
    const wrapper = factory({ propsData: { template, templateProps: { entity } } });
    await flushPromises();

    expect(wrapper.findComponent(CEntityTagsChips).props('entity')).toEqual(entity);
    expect(wrapper.text()).toContain('entity-tag');
  });

  it('keeps state, rich text, map, now and copy helpers as real components', async () => {
    const wrapper = factory({ propsData: {
      template: '<strong>{{state state}}</strong>{{map labels color="blue"}}{{now realtime=true format="HH:mm:ss"}}{{#copy value}}copy{{/copy}}',
      context: { state: 2, labels: 'map[first:one second:two]', value: 'copied' },
    } });
    await flushPromises();

    expect(wrapper.findComponent(CAlarmChip).props('value')).toBe('2');
    expect(wrapper.findAll('.v-chip')).toHaveLength(2);
    expect(wrapper.findComponent(CNow).props('format')).toBe('HH:mm:ss');
    expect(wrapper.findComponent(CCopyWrapper).props('value')).toBe('copied');
  });
});

describe('fixed helper data and events', () => {
  const host = {
    components,
    data: () => ({ alarm: { tag_colors: [{ value: 'tag', color: '#ff0000' }, { value: 'other', color: '#00ff00' }] }, selectedTags: ['tag'] }),
    computed: {
      template() {
        return sanitizeVueTemplate(`<div>${alarmOrEntityTagsHelper.call(this, { hash: {} })}</div>`);
      },
    },
    template: '<c-runtime-template :template="template" @select:tag="$emit(\'select\', $event)" @remove:tag="$emit(\'close\', $event)"/>',
  };
  const factory = generateRenderer(host, { localVue, parentComponent });

  it('passes tag data and forwards select and close events', async () => {
    const wrapper = factory();
    const tags = wrapper.findComponent(CAlarmTagsChips);

    expect(tags.props('alarm')).toEqual(wrapper.vm.alarm);
    expect(tags.props('selectedTags')).toEqual(['tag']);
    await tags.findAll('.c-alarm-tags-chips__chip').at(1).trigger('click');
    await tags.find('.c-alarm-tags-chips__chip .v-icon').trigger('click');
    await Vue.nextTick();
    expect(wrapper).toEmit('select', 'other');
    expect(wrapper).toEmit('close', 'tag');
  });
});

describe('links category remains a resolved HBS string', () => {
  const linksMixin = handlebarsLinksHelperCreator('links', 'unused');
  const factory = generateRenderer(CCompiledTemplate, { localVue, parentComponent });

  let consoleErrorSpy;

  beforeEach(() => linksMixin.beforeCreate());
  afterEach(() => {
    linksMixin.destroyed();
    consoleErrorSpy?.mockRestore();
    consoleErrorSpy = undefined;
  });

  it.each(['test', '\' " <b> &', expression])('escapes category %s and updates its context', async (category) => {
    const wrapper = factory({ propsData: {
      template: '{{links category=entity.category}}',
      context: { entity: { category } },
      templateProps: { links: {}, hasAccessForLinks: true },
    } });
    await flushPromises();

    expect(wrapper.findComponent(CLinksList).props('category')).toBe(category);
    await wrapper.setProps({ context: { entity: { category: 'updated' } } });
    await flushPromises();
    expect(wrapper.findComponent(CLinksList).props('category')).toBe('updated');
    expect(window.templateXss).toBeUndefined();
  });

  it('preserves a literal category containing braces and checks the expected Vue 2 warning', async () => {
    consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    const wrapper = factory({ propsData: {
      template: '{{links category=entity.category}}',
      context: { entity: { category: '{{literal}}' } },
      templateProps: { links: {}, hasAccessForLinks: true },
    } });
    await flushPromises();

    expect(wrapper.findComponent(CLinksList).props('category')).toBe('{{literal}}');
    expect(consoleErrorSpy).toHaveBeenCalledTimes(1);
    expect(consoleErrorSpy).toHaveBeenCalledWith(expect.stringContaining(
      'category="{{literal}}": Interpolation inside attributes has been removed.',
    ));
    await wrapper.setProps({ context: { entity: { category: 'updated' } } });
    await flushPromises();
    expect(wrapper.findComponent(CLinksList).props('category')).toBe('updated');
    expect(consoleErrorSpy).toHaveBeenCalledTimes(1);
    expect(window.templateXss).toBeUndefined();
  });

  it.each(['{{links category="test"}}', '{{links}}', '{{links category=false}}'])('preserves literals and the absent category default: %s', async (template) => {
    const wrapper = factory({ propsData: {
      template, templateProps: { links: {}, hasAccessForLinks: true },
    } });
    await flushPromises();

    expect(wrapper.findComponent(CLinksList).props('category')).toBe(template.includes('"test"') ? 'test' : null);
  });
});

describe('entities helper', () => {
  let consoleWarnSpy;
  beforeEach(() => {
    consoleWarnSpy = jest.spyOn(console, 'warn');
  });
  afterEach(() => consoleWarnSpy.mockRestore());

  const factory = generateRenderer(ServiceTemplate, { localVue, parentComponent });

  it('keeps service data, pagination, refresh and action forwarding', async () => {
    const wrapper = factory({ propsData: {
      service: { _id: 'service' },
      serviceEntities: [],
      options: { page: 1, itemsPerPage: 10 },
      totalItems: 30,
      widgetParameters: { modalTemplate: '{{entities}}' },
    } });
    await flushPromises();
    const entities = wrapper.findComponent(ServiceEntitiesList);

    expect(entities.props('service')).toEqual({ _id: 'service' });
    expect(entities.props('totalItems')).toBe(30);
    const pagination = wrapper.findComponent(CTablePagination);
    await pagination.triggerCustomEvent('update:page', 2);
    await entities.triggerCustomEvent('refresh', true);
    entities.vm.addAction({ type: 'test' });
    await Vue.nextTick();
    expect(wrapper).toEmit('update:options', { page: 2, itemsPerPage: 10 });
    expect(wrapper).toEmit('refresh', true);
    expect(wrapper).toEmit('add:action', { type: 'test' });
    expect(consoleWarnSpy).not.toHaveBeenCalled();
  });
});

describe('asynchronous request helper', () => {
  const axiosMock = new AxiosMockAdapter(axios);
  const factory = generateRenderer(CCompiledTemplate, { localVue, parentComponent });

  afterEach(() => axiosMock.reset());

  it('sanitizes HBS block output from the response before Vue compilation', async () => {
    axiosMock.onGet('/xss-regression').reply(200, { comment: `<span v-if="${expression}">response</span>{{ ${expression} }}` });
    const wrapper = factory({ propsData: {
      template: '{{#request url="/xss-regression"}}{{{comment}}}{{/request}}',
    } });
    await flushPromises();

    expect(wrapper.text()).toBe(`response{{ ${expression} }}`);
    expect(wrapper.find('[v-if]').exists()).toBe(false);
    expect(window.templateXss).toBeUndefined();
  });
});

describe('SVG registry and previews', () => {
  let consoleWarnSpy;
  beforeEach(() => {
    consoleWarnSpy = jest.spyOn(console, 'warn');
  });

  const factory = generateRenderer({
    components: { ...components, CustomIcon, IconFileSelector },
    mixins: [vuetifyCustomIconsBaseMixin],
    data: () => ({ template: sanitizeVueTemplate('<div><v-icon>$vuetify.icons.xssRegression</v-icon></div>') }),
    created() {
      this.registerIconInVuetify('xssRegression', maliciousSvg, { fillBorder: true });
    },
    beforeDestroy() {
      this.unregisterIconFromVuetify('xssRegression');
    },
    template: '<div><c-runtime-template :template="template"/><custom-icon :content="content"/><icon-file-selector ref="selector"/></div>',
    computed: { content: () => maliciousSvg },
  }, { localVue, parentComponent });

  afterEach(() => {
    delete window.templateXss;
    consoleWarnSpy.mockRestore();
  });

  it('uses native sanitized SVG in v-icon, table/edit preview and file selection', async () => {
    const wrapper = factory();
    const icon = wrapper.find('svg');

    expect(icon.attributes('viewBox')).toBe('0 0 24 24');
    expect(wrapper.find('.v-icon--fill-border').exists()).toBe(true);
    expect(wrapper.findAll('foreignObject, img, [onload], [onerror]')).toHaveLength(0);
    expect(wrapper.find('text').text()).toBe(`{{ ${expression} }}`);
    const file = new File([maliciousSvg], 'safe-marker.svg', { type: 'image/svg+xml' });
    const selector = wrapper.findComponent(IconFileSelector);
    await selector.vm.parseIcon(file);
    await flushPromises();
    const registered = selector.vm.$vuetify.icons.values[selector.vm.uniqueIconName];
    const uploaded = mount(registered.component);
    expect(uploaded.find('path').attributes('d')).toBe('M0 0');
    expect(uploaded.findAll('foreignObject, img, [onerror]')).toHaveLength(0);
    expect(uploaded.find('text').text()).toBe(`{{ ${expression} }}`);
    expect(window.templateXss).toBeUndefined();
    expect(consoleWarnSpy).not.toHaveBeenCalled();
    uploaded.destroy();
  });
});
