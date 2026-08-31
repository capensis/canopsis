<template>
  <div class="c-form-general-patterns-tabs">
    <v-tabs
      v-model="activeTab"
      :key="visibleTabsKey"
      slider-color="primary"
      align-tabs="center"
    >
      <template v-for="tab in visibleTabs" :key="`tab-head-${tab.id}`">
        <v-tab
          v-if="tab.testing"
          :value="tab.id"
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
        <v-tab
          v-else
          :value="tab.id"
          :class="tabTabClass(tab)"
          :disabled="tab['test-query'] && !!disabledTestQueryTooltip"
        >
          <template v-if="tab.general">
            {{ $t('common.general') }}
          </template>
          <template v-else-if="tab.patterns">
            {{ patternsLabel || t('common.pattern', 2) }}
          </template>
          <template v-else-if="tab.additional">
            {{ additionalLabel || $t('common.additional') }}
          </template>
          <template v-else-if="tab['test-query']">
            <v-tooltip :disabled="!disabledTestQueryTooltip" location="top">
              <template #activator="{ props: tooltipProps }">
                <span v-bind="tooltipProps">{{ $t('common.testQuery') }}</span>
              </template>
              <span>{{ disabledTestQueryTooltip }}</span>
            </v-tooltip>
          </template>
        </v-tab>
      </template>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item
        v-for="tab in visibleTabs"
        :key="`item-${tab.id}`"
        :value="tab.id"
        :class="tabItemClass(tab)"
        :eager="tab.general || tab.patterns || tab.additional"
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

        <slot
          v-else-if="tab['test-query']"
          :set-ref="setTestQueryRef"
          :form="form"
          name="test-query"
        />

        <template-testing-test-variables
          v-else-if="tab.testing"
          :ref="setTestingRef"
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

import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';
import { useI18n } from '@/hooks/i18n';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';
import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';

import { useTestVariablesTabData } from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

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
const createTab = id => ({
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
  emits: ['update:form'],
  setup(props, { emit }) {
    const slots = useSlots();
    const { t } = useI18n();

    const activeTab = ref(FORM_GENERAL_PATTERNS_TABS.general);

    const generalElement = ref(null);
    const patternsElement = ref(null);
    const additionalElement = ref(null);
    const testQueryElement = ref(null);
    const testingElement = ref(null);

    const setGeneralRef = refElement => (generalElement.value = refElement);
    const setPatternsRef = refElement => (patternsElement.value = refElement);
    const setAdditionalRef = refElement => (additionalElement.value = refElement);
    const setTestQueryRef = refElement => (testQueryElement.value = refElement);
    const setTestingRef = refElement => (testingElement.value = refElement);

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasPatternsError } = useValidationElementChildren(patternsElement);
    const { hasChildrenError: hasAdditionalError } =
      useValidationElementChildren(additionalElement);
    const { hasChildrenError: hasTestQueryError } =
      useValidationElementChildren(testQueryElement);
    const { hasChildrenError: hasTestingError } = useValidationElementChildren(testingElement);
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
    const visibleTabsKey = computed(() => visibleTabs.value.map(tab => tab.id).join('|'));

    const patternsTabIndex = computed(
      () => visibleTabs.value.some(tab => tab.id === FORM_GENERAL_PATTERNS_TABS.patterns)
        ? FORM_GENERAL_PATTERNS_TABS.patterns
        : FORM_GENERAL_PATTERNS_TABS.general,
    );

    const testingTabIndex = computed(() => visibleTabs.value.some(tab => tab.id === FORM_GENERAL_PATTERNS_TABS.testing)
        ? FORM_GENERAL_PATTERNS_TABS.testing
        : null,
    );

    const isActiveTestingTab = computed(
      () => testingTabIndex.value !== null && activeTab.value === testingTabIndex.value,
    );

    const tabHasError = tab => (tab.general && hasGeneralError.value) ||
      (tab.patterns && hasPatternsError.value) ||
      (tab.additional && hasAdditionalError.value) ||
      (tab['test-query'] && hasTestQueryError.value) ||
      (tab.testing && hasTestingError.value);

    const tabTabClass = tab => ({
      'text-error': tabHasError(tab),
      'v-tab--tooltip': tab['test-query'],
    });

    const tabItemClass = tab => (tab.general || tab.patterns || tab.additional ? 'pt-4' : '');

    const expandNeededTab = computed(() => {
      if (props.aiChatExpandTabKeys) {
        return Object.fromEntries(
          Object.entries(props.aiChatExpandTabKeys).map(([fieldKey, tabId]) => [
            fieldKey,
            visibleTabs.value.some(tab => tab.id === tabId)
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
      if (!tabs.some(tab => tab.id === activeTab.value)) {
        activeTab.value = tabs[0]?.id ?? FORM_GENERAL_PATTERNS_TABS.general;
      }
    });

    onMounted(() => {
      if (isNumber(props.type)) {
        fetchTemplateVarsList();
      }

      if (hasTemplateTestingTab.value) {
        fetchCopyVarsList();
      }
    });

    return {
      activeTab,
      t,

      visibleTabs,
      visibleTabsKey,

      hasGeneralError,
      hasPatternsError,
      hasAdditionalError,

      setGeneralRef,
      setPatternsRef,
      setAdditionalRef,
      setTestQueryRef,
      setTestingRef,

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
