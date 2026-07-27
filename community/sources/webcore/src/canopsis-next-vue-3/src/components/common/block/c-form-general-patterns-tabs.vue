<template>
  <div class="c-form-general-patterns-tabs">
    <v-tabs v-model="activeTab" :key="visibleTabsKey" slider-color="primary" align-tabs="center">
      <v-tab
        v-if="hasPatternsTab && reverse"
        :value="FORM_GENERAL_PATTERNS_TABS.patterns"
        :class="tabTabClass(tabsByIds[FORM_GENERAL_PATTERNS_TABS.patterns])"
      >
        {{ patternsLabel || t('common.pattern', 2) }}
      </v-tab>

      <v-tab
        v-if="hasGeneralTab"
        :value="FORM_GENERAL_PATTERNS_TABS.general"
        :class="tabTabClass(tabsByIds[FORM_GENERAL_PATTERNS_TABS.general])"
      >
        {{ $t('common.general') }}
      </v-tab>

      <v-tab
        v-if="hasPatternsTab && !reverse"
        :value="FORM_GENERAL_PATTERNS_TABS.patterns"
        :class="tabTabClass(tabsByIds[FORM_GENERAL_PATTERNS_TABS.patterns])"
      >
        {{ patternsLabel || t('common.pattern', 2) }}
      </v-tab>

      <v-tab
        v-if="hasAdditionalTab"
        :value="FORM_GENERAL_PATTERNS_TABS.additional"
        :class="tabTabClass(tabsByIds[FORM_GENERAL_PATTERNS_TABS.additional])"
      >
        {{ additionalLabel || $t('common.additional') }}
      </v-tab>

      <v-tab
        v-if="hasTestQueryTab"
        :value="FORM_GENERAL_PATTERNS_TABS.testQuery"
        :class="tabTabClass(tabsByIds[FORM_GENERAL_PATTERNS_TABS.testQuery])"
        :disabled="!!disabledTestQueryTooltip"
      >
        <v-tooltip :disabled="!disabledTestQueryTooltip" location="top">
          <template #activator="{ props: tooltipProps }">
            <span v-bind="tooltipProps">{{ $t('common.testQuery') }}</span>
          </template>
          <span>{{ disabledTestQueryTooltip }}</span>
        </v-tooltip>
      </v-tab>

      <v-tab
        v-if="hasTemplateTestingTab"
        :value="FORM_GENERAL_PATTERNS_TABS.testing"
        :disabled="isEmptyVariablesFields"
        class="v-tab--tooltip"
      >
        <v-tooltip :disabled="!isEmptyVariablesFields" location="top">
          <template #activator="{ props: tooltipProps }">
            <span v-bind="tooltipProps" class="test-variables-tab__wrapper">
              {{ t('templateTesting.testVariables') }}
            </span>
          </template>

          <span>{{ t('templateTesting.testVariablesDisabledTooltip') }}</span>
        </v-tooltip>
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item
        v-for="tab in visibleTabs"
        :key="`item-${tab.id}`"
        :value="tab.id"
        :class="tabItemClass(tab)"
        :disabled="tab.testing && isEmptyVariablesFields"
      >
        <slot
          v-if="tab.general"
          :set-ref="setGeneralRef"
          :copy-vars="copyVars"
          :template-vars="templateVars"
          :pending="templateTestingPending"
          name="general"
        />

        <slot
          v-else-if="tab.additional"
          :set-ref="setAdditionalRef"
          :template-vars="templateVars"
          :copy-vars="copyVars"
          :pending="templateTestingPending"
          name="additional"
        />

        <slot
          v-else-if="tab.patterns"
          :set-ref="setPatternsRef"
          :template-vars="templateVars"
          :copy-vars="copyVars"
          :pending="templateTestingPending"
          name="patterns"
        />

        <slot v-else-if="tab['test-query']" :form="form" name="test-query" />

        <template-testing-test-variables
          v-else-if="tab.testing"
          :general-form="form"
          :variables-fields="variablesFields"
          :template-vars="templateVars"
          :rule-id="ruleId"
          :type="type"
          :active="isActiveTestingTab"
        />
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { isNumber } from 'lodash';
import { computed, ref, toRef, useSlots, watch, onMounted } from 'vue';

import { FORM_GENERAL_PATTERNS_TABS } from '@/constants';

import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import { useTestVariablesTabData } from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

import { useI18n } from '@/hooks/i18n';

import TemplateTestingTestVariables from '@/components/other/template-testing/test-variables/template-testing-test-variables.vue';

/**
 * Builds a tab descriptor for `c-form-general-patterns-tabs` with a single active type flag.
 *
 * @param {string} id - Tab id from `FORM_GENERAL_PATTERNS_TABS`.
 * @returns {{
 *   id: string,
 *   general: boolean,
 *   patterns: boolean,
 *   'test-query': boolean,
 *   testing: boolean,
 *   additional: boolean,
 * }}
 */
const createTab = (id) => ({
  id,
  [FORM_GENERAL_PATTERNS_TABS.general]: id === FORM_GENERAL_PATTERNS_TABS.general,
  [FORM_GENERAL_PATTERNS_TABS.patterns]: id === FORM_GENERAL_PATTERNS_TABS.patterns,
  [FORM_GENERAL_PATTERNS_TABS.testQuery]: id === FORM_GENERAL_PATTERNS_TABS.testQuery,
  [FORM_GENERAL_PATTERNS_TABS.testing]: id === FORM_GENERAL_PATTERNS_TABS.testing,
  [FORM_GENERAL_PATTERNS_TABS.additional]: id === FORM_GENERAL_PATTERNS_TABS.additional,
});

export default {
  components: { TemplateTestingTestVariables },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    ruleId: {
      type: String,
      required: false,
    },
    type: {
      type: Number,
      required: false,
    },
    reverse: {
      type: Boolean,
      default: false,
    },
    patternsLabel: {
      type: String,
      default: '',
    },
    additionalLabel: {
      type: String,
      default: '',
    },
    disabledTestQueryTooltip: {
      type: String,
      required: false,
      default: '',
    },
    hideGeneral: {
      type: Boolean,
      default: false,
    },
    withTemplateTesting: {
      type: Boolean,
      default: true,
    },
    aiChatExpandTabKeys: {
      type: Object,
      default: null,
    },
  },
  setup(props, { emit }) {
    const slots = useSlots();
    const { t } = useI18n();

    const activeTab = ref(FORM_GENERAL_PATTERNS_TABS.general);

    const generalElement = ref(null);
    const patternsElement = ref(null);
    const additionalElement = ref(null);

    const setGeneralRef = (refElement) => (generalElement.value = refElement);
    const setPatternsRef = (refElement) => (patternsElement.value = refElement);
    const setAdditionalRef = (refElement) => (additionalElement.value = refElement);

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasPatternsError } = useValidationElementChildren(patternsElement);
    const { hasChildrenError: hasAdditionalError } =
      useValidationElementChildren(additionalElement);
    const {
      vars: copyVars,
      pending: copyVarsPending,
      fetchList: fetchCopyVarsList,
    } = useCopyVarsList({
      type: toRef(props, 'type'),
    });

    const {
      vars: templateVars,
      pending: templateVarsPending,
      fetchList: fetchTemplateVarsList,
    } = useTemplateVarsList({
      type: toRef(props, 'type'),
      form: toRef(props, 'form'),
    });

    const {
      hasAccess: hasAccessToTemplateTesting,
      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesTabData(props, toRef(props, 'type'), emit);

    const hasTemplateTestingTab = computed(
      () => props.withTemplateTesting && isNumber(props.type) && hasAccessToTemplateTesting.value,
    );
    const templateTestingPending = computed(
      () => templateVarsPending.value || copyVarsPending.value,
    );

    const hasGeneralSlot = computed(() => Boolean(slots.general));
    const hasPatternsSlot = computed(() => Boolean(slots.patterns));
    const hasTestQuerySlot = computed(() => Boolean(slots['test-query']));
    const hasAdditionalSlot = computed(() => Boolean(slots.additional));
    const hasGeneralTab = computed(() => hasGeneralSlot.value && !props.hideGeneral);
    const hasPatternsTab = computed(() => hasPatternsSlot.value);
    const hasTestQueryTab = computed(() => hasTestQuerySlot.value);
    const hasAdditionalTab = computed(() => hasAdditionalSlot.value);

    const tabsByIds = computed(() =>
      Object.fromEntries(
        Object.values(FORM_GENERAL_PATTERNS_TABS).map((tabId) => [tabId, createTab(tabId)]),
      ),
    );

    const visibleTabs = computed(() => {
      const generalTab =
        hasGeneralSlot.value && !props.hideGeneral && createTab(FORM_GENERAL_PATTERNS_TABS.general);
      const patternsTab = hasPatternsSlot.value && createTab(FORM_GENERAL_PATTERNS_TABS.patterns);
      const testQueryTab =
        hasTestQuerySlot.value && createTab(FORM_GENERAL_PATTERNS_TABS.testQuery);
      const testingTab =
        hasTemplateTestingTab.value && createTab(FORM_GENERAL_PATTERNS_TABS.testing);
      const additionalTab =
        hasAdditionalSlot.value && createTab(FORM_GENERAL_PATTERNS_TABS.additional);

      const tabs = props.reverse
        ? [patternsTab, generalTab, additionalTab]
        : [generalTab, additionalTab, patternsTab];

      tabs.push(testQueryTab, testingTab);

      return tabs.filter(Boolean);
    });
    const visibleTabsKey = computed(() => visibleTabs.value.map((tab) => tab.id).join('|'));

    const patternsTabIndex = computed(() =>
      visibleTabs.value.some((tab) => tab.id === FORM_GENERAL_PATTERNS_TABS.patterns)
        ? FORM_GENERAL_PATTERNS_TABS.patterns
        : FORM_GENERAL_PATTERNS_TABS.general,
    );

    const testingTabIndex = computed(() =>
      visibleTabs.value.some((tab) => tab.id === FORM_GENERAL_PATTERNS_TABS.testing)
        ? FORM_GENERAL_PATTERNS_TABS.testing
        : null,
    );

    const isActiveTestingTab = computed(
      () => testingTabIndex.value !== null && activeTab.value === testingTabIndex.value,
    );

    const tabTabClass = (tab) => ({
      'error--text':
        (tab.general && hasGeneralError.value) ||
        (tab.patterns && hasPatternsError.value) ||
        (tab.additional && hasAdditionalError.value),
      'v-tab--tooltip': tab['test-query'],
    });

    const tabItemClass = (tab) => (tab.general || tab.patterns || tab.additional ? 'pt-4' : '');

    const expandNeededTab = computed(() => {
      if (props.aiChatExpandTabKeys) {
        return Object.fromEntries(
          Object.entries(props.aiChatExpandTabKeys).map(([fieldKey, tabId]) => [
            fieldKey,
            visibleTabs.value.some((tab) => tab.id === tabId)
              ? tabId
              : FORM_GENERAL_PATTERNS_TABS.general,
          ]),
        );
      }

      return patternsTabIndex.value;
    });

    useAiChatExpand({
      activeTab,
      neededTab: expandNeededTab,
    });

    watch(
      () => props.hideGeneral,
      (value) => {
        if (value) {
          activeTab.value = patternsTabIndex.value;
        }
      },
    );

    watch(visibleTabs, (tabs) => {
      if (!tabs.some((tab) => tab.id === activeTab.value)) {
        activeTab.value = tabs[0]?.id ?? FORM_GENERAL_PATTERNS_TABS.general;
      }
    });

    onMounted(() => {
      if (hasTemplateTestingTab.value) {
        fetchCopyVarsList();
        fetchTemplateVarsList();
      }
    });

    return {
      activeTab,
      t,
      FORM_GENERAL_PATTERNS_TABS,

      visibleTabs,
      visibleTabsKey,
      tabsByIds,
      hasGeneralTab,
      hasPatternsTab,
      hasTestQueryTab,
      hasAdditionalTab,

      hasGeneralError,
      hasPatternsError,

      setGeneralRef,
      setPatternsRef,
      setAdditionalRef,

      copyVars,
      templateVars,
      hasTemplateTestingTab,
      variablesFields,
      isEmptyVariablesFields,
      isActiveTestingTab,
      templateTestingPending,

      tabTabClass,
      tabItemClass,
    };
  },
};
</script>

<style lang="scss" scoped>
.test-variables-tab__wrapper {
  max-width: 360px;
}
</style>
