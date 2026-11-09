import { flushPromises, mount } from '@vue/test-utils';
import axios from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';
import { createPinia } from 'pinia';
import { createMemoryHistory, createRouter, RouterLink } from 'vue-router';

import { vuetify } from '@/plugins/vuetify';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';
import { alarmOrEntityTagsHelper } from '@/helpers/handlebars/tags-helper';
import { sanitizeHtml } from '@/helpers/html';

import CAlarmActionsChips from '@/components/common/chips/c-alarm-actions-chips.vue';
import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';
import CAlarmTagsChips from '@/components/common/chips/c-alarm-tags-chips.vue';
import CChip from '@/components/common/chips/c-chip.vue';
import CIcon from '@/components/common/icon/c-icon.vue';
import CLayout from '@/components/common/layout/c-layout.vue';
import CPagination from '@/components/common/pagination/c-pagination.vue';
import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CCopyWrapper from '@/components/common/runtime-template/c-copy-wrapper.vue';
import CNow from '@/components/common/runtime-template/c-now.vue';
import CRequestHelper from '@/components/common/runtime-template/c-request-helper.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CSelect from '@/components/common/vuetify/c-select.vue';
import CItemsPerPageField from '@/components/forms/fields/c-items-per-page-field.vue';
import ServiceEntitiesList from '@/components/other/service/partials/service-entities-list.vue';
import ServiceTemplate from '@/components/other/service/partials/service-template.vue';

const createGlobalOptions = () => {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/:pathMatch(.*)*', component: { template: '<div />' } }],
  });

  return {
    plugins: [createPinia(), vuetify, router],
    components: {
      CAlarmActionsChips,
      CAlarmChip,
      CAlarmTagsChips,
      CChip,
      CLayout,
      CIcon,
      CPagination,
      CTablePagination,
      CSelect,
      CItemsPerPageField,
      CCompiledTemplate,
      CCopyWrapper,
      CNow,
      CRequestHelper,
      CRuntimeTemplate,
    },
    provide: { $system: { theme: {}, dark: false } },
    mocks: { $popups: { success: vi.fn(), error: vi.fn() } },
  };
};

const factory = (template, options = {}) => mount(CCompiledTemplate, {
  props: { template, ...options },
  global: createGlobalOptions(),
});

describe('real Vue components in Handlebars templates', () => {
  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it('preserves Handlebars conditions in table comments through editing and rendering', async () => {
    const template = sanitizeHtml('<table><tbody><!--{{#each items}}--><!--{{#if enabled}}--><tr><td>{{name}}</td><td>{{state state}}</td></tr><!--{{/if}}--><!--{{/each}}--></tbody></table>');
    const wrapper = factory(template, {
      context: {
        items: [
          { enabled: true, name: 'Shown', state: 1 },
          { enabled: false, name: 'Hidden', state: 2 },
        ],
      },
    });

    await flushPromises();

    expect(wrapper.findAll('tbody > tr')).toHaveLength(1);
    expect(wrapper.find('td').text()).toBe('Shown');
    expect(wrapper.findComponent(CAlarmChip).props('value')).toBe('1');

    await wrapper.setProps({
      context: { items: [{ enabled: true, name: 'Updated', state: 3 }] },
    });
    await flushPromises();

    expect(wrapper.findAll('tbody > tr')).toHaveLength(1);
    expect(wrapper.find('td').text()).toBe('Updated');
    expect(wrapper.findComponent(CAlarmChip).props('value')).toBe('3');
    wrapper.unmount();
  });

  it.each([true, false])('evaluates a condition in comments around a whole table: %s', async (enabled) => {
    const template = sanitizeHtml('<!--{{#if enabled}}--><table><tbody><tr><td>yes</td></tr></tbody></table><!--{{else}}--><p>no</p><!--{{/if}}-->');
    const wrapper = factory(template, { context: { enabled } });

    await flushPromises();

    expect(wrapper.find('table').exists()).toBe(enabled);
    expect(wrapper.text()).toBe(enabled ? 'yes' : 'no');
    wrapper.unmount();
  });

  it('renders the state helper as an actual alarm chip', async () => {
    const wrapper = factory('<p>{{state alarm.v.state.val}}</p>', {
      context: { alarm: { v: { state: { val: 2 } } } },
    });

    await flushPromises();

    expect(wrapper.findComponent(CAlarmChip).props('value')).toBe('2');
    expect(wrapper.find('.chip').exists()).toBe(true);
    expect(wrapper.find('c-alarm-chip').exists()).toBe(false);
    wrapper.unmount();
  });

  it('renders legacy icons through c-icon with size compatibility and overflow clipping', async () => {
    const wrapper = factory('<v-icon small color="red">home</v-icon>');

    await flushPromises();

    expect(wrapper.findComponent(CIcon).props('small')).toBe(true);
    expect(wrapper.find('i').classes()).toContain('v-icon--size-small');
    expect(wrapper.find('i').classes()).toContain('overflow-hidden');
    expect(wrapper.find('i').classes()).toContain('text-red');
    expect(wrapper.find('i').text()).toBe('home');
    wrapper.unmount();
  });

  it('renders the map helper as Vuetify rows and chips', async () => {
    const wrapper = factory('{{map labels color="blue"}}', {
      context: { labels: 'map[first:one second:two]' },
    });

    await flushPromises();

    expect(wrapper.find('.v-row').exists()).toBe(true);
    expect(wrapper.findAll('.v-chip').map(chip => chip.text())).toEqual(['first:one', 'second:two']);
    wrapper.unmount();
  });

  it('renders internal links as RouterLink instances with functional navigation', async () => {
    const wrapper = factory('{{internal-link href="/alarms" text="Alarms"}}');

    await flushPromises();

    expect(wrapper.findComponent(RouterLink).props('to')).toBe('/alarms');
    expect(wrapper.find('a').attributes('href')).toBe('/alarms');
    await wrapper.find('a').trigger('click');
    await flushPromises();
    expect(wrapper.vm.$router.currentRoute.value.path).toBe('/alarms');
    wrapper.unmount();
  });

  it('renders the copy helper and retains its click behavior', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);

    vi.stubGlobal('navigator', { clipboard: { writeText } });

    const wrapper = factory('{{#copy value}}<strong>Copy value</strong>{{/copy}}', {
      context: { value: 'copied text' },
    });

    await flushPromises();
    expect(wrapper.findComponent(CCopyWrapper).props('value')).toBe('copied text');
    await wrapper.find('strong').trigger('click');
    await flushPromises();
    expect(writeText).toHaveBeenCalledWith('copied text');
    wrapper.unmount();
  });

  it('renders the realtime now helper and updates it on subsequent ticks', async () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-11-03T12:00:00Z'));

    const wrapper = factory('{{now format="HH:mm:ss" timezone="UTC" realtime=true}}');

    await vi.advanceTimersByTimeAsync(0);
    expect(wrapper.findComponent(CNow).exists()).toBe(true);
    expect(wrapper.vm.compiledTemplate).toContain('format="HH:mm:ss"');
    expect(wrapper.findComponent(CNow).props('format')).toBe('HH:mm:ss');
    expect(wrapper.text()).toBe('12:00:00');
    await vi.advanceTimersByTimeAsync(1000);
    expect(wrapper.text()).toBe('12:00:01');
    wrapper.unmount();
  });

  it('renders real alarm tags, updates their data and forwards selection and close events', async () => {
    registerHelper('tags', alarmOrEntityTagsHelper);

    const select = vi.fn();
    const close = vi.fn();
    const alarm = { tag_colors: [{ value: 'first', color: '#112233' }] };
    const wrapper = factory('{{tags}}', {
      context: { alarm },
      templateProps: { alarm, selectedTags: [] },
      'onSelect:tag': select,
      'onRemove:tag': close,
    });

    await flushPromises();
    expect(wrapper.findComponent(CAlarmTagsChips).props('alarm')).toEqual(alarm);
    await wrapper.find('.c-chip').trigger('click');
    expect(select).toHaveBeenCalledWith('first');

    await wrapper.setProps({ templateProps: { alarm, selectedTags: ['first'] } });
    await wrapper.find('.c-chip .v-icon').trigger('click');
    expect(close).toHaveBeenCalledWith('first');
    wrapper.unmount();
    unregisterHelper('tags');
  });


  it('renders the actual service list inside a composed template and forwards pagination events', async () => {
    const service = { _id: 'service-id' };
    const options = { page: 1, itemsPerPage: 10 };
    const wrapper = mount(ServiceTemplate, {
      props: {
        service,
        options,
        serviceEntities: [],
        totalItems: 0,
        widgetParameters: { modalTemplate: '<h3>Service</h3>{{entities name="entity._id"}}' },
      },
      global: createGlobalOptions(),
    });

    await flushPromises();

    const list = wrapper.findComponent(ServiceEntitiesList);

    expect(wrapper.findComponent(CCompiledTemplate).exists()).toBe(true);
    expect(wrapper.find('h3').text()).toBe('Service');
    expect(list.props()).toMatchObject({ service, options, entityNameField: 'entity._id' });
    list.vm.updatePage(2);
    await flushPromises();
    expect(wrapper.emitted('update:options')).toEqual([[{ page: 2, itemsPerPage: 10 }]]);
    wrapper.unmount();
  });

  it('renders components from the asynchronous request helper response', async () => {
    const adapter = new AxiosMockAdapter(axios);

    adapter.onGet('/template-data').reply(200, { message: 'Loaded', state: 1 });

    const wrapper = factory('{{#request url="/template-data"}}<strong>{{message}}</strong>{{state state}}{{/request}}');

    await flushPromises();
    expect(wrapper.findComponent(CRequestHelper).exists()).toBe(true);
    expect(wrapper.find('strong').text()).toBe('Loaded');
    expect(wrapper.findComponent(CAlarmChip).props('value')).toBe('1');
    wrapper.unmount();
    adapter.restore();
  });
});
