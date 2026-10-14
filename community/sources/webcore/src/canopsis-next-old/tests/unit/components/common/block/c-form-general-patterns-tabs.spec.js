import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import {
  createMockedStoreModules,
  createAuthModule,
  createCopyVarsModule,
  createTemplateVarsModule,
} from '@unit/utils/store';

import { FORM_GENERAL_PATTERNS_TABS, TEMPLATE_TESTING_TEST_TYPES, USER_PERMISSIONS } from '@/constants';

import CFormGeneralPatternsTabs from '@/components/common/block/c-form-general-patterns-tabs.vue';

const { authModule, currentUserPermissionsById } = createAuthModule();
const {
  copyVarsModule,
  fetchEventFiltersVarsWithoutStore,
  fetchDynamicInfosVarsWithoutStore,
} = createCopyVarsModule();
const {
  templateVarsModule,
  fetchDeclareTicketRulesVarsWithoutStore,
} = createTemplateVarsModule();

const defaultSlots = {
  general: '<div class="general-slot" />',
  patterns: '<div class="patterns-slot" />',
};

const stubs = {
  'template-testing-test-variables': true,
  'template-testing-test-variables-tab': true,
};

const selectTabItems = wrapper => wrapper.findAll('v-tab-stub');
const selectGeneralTab = wrapper => selectTabItems(wrapper).wrappers.find(tab => tab.text() === 'General');
const selectPatternsTab = wrapper => selectTabItems(wrapper).wrappers.find(tab => /pattern/i.test(tab.text()));

const createStore = ({ withTemplateTestingAccess = false } = {}) => {
  currentUserPermissionsById.mockReturnValue(
    withTemplateTestingAccess
      ? { [USER_PERMISSIONS.technical.templateTesting]: { actions: [] } }
      : {},
  );

  return createMockedStoreModules([authModule, copyVarsModule, templateVarsModule]);
};

describe('c-form-general-patterns-tabs', () => {
  const factory = generateShallowRenderer(CFormGeneralPatternsTabs, {
    stubs,
    store: createStore(),
  });
  const snapshotFactory = generateRenderer(CFormGeneralPatternsTabs, {
    stubs,
    store: createStore(),
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
    expect(tabs.at(1).text()).toBe('General');
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
    expect(wrapper.vm.activeTab).toBe(0);
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
    expect(wrapper.vm.activeTab).toBe(0);
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

    const testQueryTab = selectTabItems(wrapper).wrappers.find(tab => tab.props('disabled'));

    expect(wrapper.vm.visibleTabs.some(tab => tab['test-query'])).toBe(true);
    expect(testQueryTab).toBeDefined();
    expect(testQueryTab.props('disabled')).toBe(true);
  });

  test('Template vars lists are fetched on mount when template testing tab is available', async () => {
    factory({
      store: createStore({ withTemplateTestingAccess: true }),
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
        form: { name: 'rule-name' },
      },
      slots: defaultSlots,
    });

    await flushPromises();

    expect(fetchDeclareTicketRulesVarsWithoutStore).toHaveBeenCalledTimes(1);
    expect(fetchEventFiltersVarsWithoutStore).not.toHaveBeenCalled();
    expect(fetchDynamicInfosVarsWithoutStore).not.toHaveBeenCalled();
  });

  test('Template vars lists are not fetched on mount when template testing tab is unavailable', async () => {
    factory({
      store: createStore({ withTemplateTestingAccess: false }),
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
        form: { name: 'rule-name' },
      },
      slots: defaultSlots,
    });

    await flushPromises();

    expect(fetchDeclareTicketRulesVarsWithoutStore).not.toHaveBeenCalled();
    expect(fetchEventFiltersVarsWithoutStore).not.toHaveBeenCalled();
    expect(fetchDynamicInfosVarsWithoutStore).not.toHaveBeenCalled();
  });

  test('Testing tab is rendered when user has access and type is defined', () => {
    const wrapper = factory({
      store: createStore({ withTemplateTestingAccess: true }),
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
      },
      slots: defaultSlots,
    });

    expect(wrapper.find('template-testing-test-variables-tab-stub').exists()).toBe(true);
    expect(wrapper.find('template-testing-test-variables-stub').exists()).toBe(true);
  });

  test('Visible tabs contain expected tab ids', () => {
    const wrapper = factory({
      store: createStore({ withTemplateTestingAccess: true }),
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

  test('Renders `c-form-general-patterns-tabs` with default slots', async () => {
    const wrapper = snapshotFactory({
      slots: defaultSlots,
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-form-general-patterns-tabs` with reverse prop', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        reverse: true,
      },
      slots: defaultSlots,
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-form-general-patterns-tabs` with all optional slots', async () => {
    const wrapper = snapshotFactory({
      store: createStore({ withTemplateTestingAccess: true }),
      propsData: {
        type: TEMPLATE_TESTING_TEST_TYPES.declareTicketRule,
        patternsLabel: 'Custom patterns',
        additionalLabel: 'Custom additional',
        form: { name: 'rule-name' },
        ruleId: 'rule-id',
      },
      slots: {
        ...defaultSlots,
        additional: '<div class="additional-slot" />',
        'test-query': '<div class="test-query-slot" />',
      },
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-form-general-patterns-tabs` with hidden general tab', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        hideGeneral: true,
      },
      slots: defaultSlots,
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });
});
