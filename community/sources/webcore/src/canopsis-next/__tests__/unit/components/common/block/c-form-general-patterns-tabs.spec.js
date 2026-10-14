import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import {
  DYNAMIC_INFO_INFORMATION_TYPES,
  FORM_GENERAL_PATTERNS_TABS,
  TEMPLATE_TESTING_TEST_TYPES,
} from '@/constants';

import CFormGeneralPatternsTabs from '@/components/common/block/c-form-general-patterns-tabs.vue';

const hooks = vi.hoisted(() => ({
  hasAccess: { value: false, __v_isRef: true },
  fetchCopyVarsList: vi.fn(),
  fetchTemplateVarsList: vi.fn(),
}));

vi.mock('@/hooks/auth', () => ({
  useCanPermission: () => ({ hasAccess: hooks.hasAccess }),
}));

vi.mock('@/hooks/vars/copy', () => ({
  useCopyVarsList: () => ({
    vars: { value: {}, __v_isRef: true },
    pending: { value: false, __v_isRef: true },
    fetchList: hooks.fetchCopyVarsList,
  }),
}));

vi.mock('@/hooks/vars/template', () => ({
  useTemplateVarsList: () => ({
    vars: { value: {}, __v_isRef: true },
    pending: { value: false, __v_isRef: true },
    fetchList: hooks.fetchTemplateVarsList,
  }),
}));

const defaultSlots = {
  general: '<div class="general-slot" />',
  patterns: '<div class="patterns-slot" />',
};

const stubs = {
  'template-testing-test-variables': true,
  'v-tabs': { template: '<div class="v-tabs"><slot /></div>' },
  'v-tab': {
    props: ['value', 'disabled'],
    template: '<button class="v-tab" :disabled="disabled"><slot /></button>',
  },
  'v-window': { template: '<div class="v-window"><slot /></div>' },
  'v-window-item': {
    props: ['value'],
    template: '<div class="v-window-item"><slot /></div>',
  },
};

const selectTabItems = wrapper => wrapper.findAll('button.v-tab');
const selectGeneralTab = wrapper =>
  selectTabItems(wrapper).find(tab => tab.text() === 'common.general');
const selectPatternsTab = wrapper =>
  selectTabItems(wrapper).find(tab => /pattern/i.test(tab.text()));

describe('c-form-general-patterns-tabs', () => {
  const factory = generateShallowRenderer(CFormGeneralPatternsTabs, {
    stubs,
  });

  beforeEach(() => {
    hooks.hasAccess.value = false;
  });

  afterEach(() => {
    hooks.fetchCopyVarsList.mockClear();
    hooks.fetchTemplateVarsList.mockClear();
  });

  test('General and patterns tabs are rendered by default', () => {
    const wrapper = factory({
      slots: defaultSlots,
    });

    expect(selectGeneralTab(wrapper).exists()).toBe(true);
    expect(selectPatternsTab(wrapper).exists()).toBe(true);
    expect(wrapper.find('.general-slot').exists()).toBe(true);
    expect(wrapper.find('.patterns-slot').exists()).toBe(true);
  });

  test('Patterns tab is rendered before general tab when reverse is true', () => {
    const wrapper = factory({
      propsData: {
        reverse: true,
      },
      slots: defaultSlots,
    });

    const tabs = selectTabItems(wrapper);

    expect(tabs.at(0).text()).toContain('Pattern');
    expect(tabs.at(1).text()).toBe('common.general');
  });

  test('General tab is hidden when hideGeneral is true', async () => {
    const wrapper = factory({
      propsData: {
        hideGeneral: true,
      },
      slots: defaultSlots,
    });

    expect(selectGeneralTab(wrapper)).toBeUndefined();
    expect(selectPatternsTab(wrapper).exists()).toBe(true);
    expect(wrapper.vm.activeTab).toBe(FORM_GENERAL_PATTERNS_TABS.general);
  });

  test('Active tab switches to patterns when hideGeneral becomes true', async () => {
    const wrapper = factory({
      propsData: {
        hideGeneral: false,
      },
      slots: defaultSlots,
    });

    wrapper.setProps({ hideGeneral: true });

    await wrapper.vm.$nextTick();

    expect(selectGeneralTab(wrapper)).toBeUndefined();
    expect(wrapper.vm.activeTab).toBe(FORM_GENERAL_PATTERNS_TABS.patterns);
  });

  test('Custom patterns and additional labels are rendered', () => {
    const wrapper = factory({
      propsData: {
        patternsLabel: 'Custom patterns',
        additionalLabel: 'Custom additional',
      },
      slots: {
        ...defaultSlots,
        additional: '<div class="additional-slot" />',
      },
    });

    const additionalTabIndex = wrapper.vm.visibleTabs.findIndex(
      tab => tab.id === FORM_GENERAL_PATTERNS_TABS.additional,
    );

    expect(selectPatternsTab(wrapper).text()).toBe('Custom patterns');
    expect(selectTabItems(wrapper).at(additionalTabIndex).text()).toBe('Custom additional');
    expect(wrapper.find('.additional-slot').exists()).toBe(true);
  });

  test('Test query tab is disabled when disabledTestQueryTooltip is set', () => {
    const wrapper = factory({
      propsData: {
        disabledTestQueryTooltip: 'Test query is disabled',
      },
      slots: {
        ...defaultSlots,
        'test-query': '<div class="test-query-slot" />',
      },
    });

    const testQueryTab = selectTabItems(wrapper).find(
      tab => tab.attributes('disabled') !== undefined,
    );

    expect(wrapper.vm.visibleTabs.some(tab => tab['test-query'])).toBe(true);
    expect(testQueryTab).toBeDefined();
    expect(testQueryTab.attributes('disabled')).toBeDefined();
  });

  test('Template vars lists are fetched on mount when template testing tab is available', async () => {
    hooks.hasAccess.value = true;
    factory({
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
        form: { name: 'rule-name' },
      },
      slots: defaultSlots,
    });

    await flushPromises();

    expect(hooks.fetchTemplateVarsList).toHaveBeenCalledTimes(1);
    expect(hooks.fetchCopyVarsList).toHaveBeenCalledTimes(1);
  });

  test('Template vars lists are not fetched on mount when template testing tab is unavailable', async () => {
    factory({
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
        form: { name: 'rule-name' },
      },
      slots: defaultSlots,
    });

    await flushPromises();

    expect(hooks.fetchTemplateVarsList).toHaveBeenCalledTimes(1);
    expect(hooks.fetchCopyVarsList).not.toHaveBeenCalled();
  });

  test('Testing tab is rendered when user has access and type is defined', () => {
    hooks.hasAccess.value = true;
    const wrapper = factory({
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
      },
      slots: defaultSlots,
    });

    expect(wrapper.find('template-testing-test-variables-stub').exists()).toBe(true);
  });

  test('Template Studio fields update the provided form context', () => {
    const form = {
      infos: [
        {
          key: 'dynamic-info-key',
          type: DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate,
          value: 'old template',
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.dynamicInfo,
        form,
      },
      slots: defaultSlots,
    });

    wrapper.vm.variablesFields[0].updateField('new template');

    expect(globalThis.__TEST_FORM__.values.value.infos[0].value).toBe('new template');
  });

  test('Visible tabs contain expected tab ids', () => {
    hooks.hasAccess.value = true;
    const wrapper = factory({
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
        reverse: true,
      },
      slots: {
        ...defaultSlots,
        additional: '<div class="additional-slot" />',
        'test-query': '<div class="test-query-slot" />',
      },
    });

    expect(wrapper.vm.visibleTabs.map(tab => tab.id)).toEqual([
      FORM_GENERAL_PATTERNS_TABS.patterns,
      FORM_GENERAL_PATTERNS_TABS.general,
      FORM_GENERAL_PATTERNS_TABS.additional,
      FORM_GENERAL_PATTERNS_TABS.testQuery,
      FORM_GENERAL_PATTERNS_TABS.testing,
    ]);
  });
});
