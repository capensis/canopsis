<template>
  <v-tabs
    v-model="activeTab"
    slider-color="primary"
    centered
  >
    <template v-for="tab in visibleTabs">
      <template-testing-test-variables-tab
        v-if="tab.testing"
        :key="`tab-head-testing-${tab.id}`"
        :disabled="isEmptyVariablesFields"
      />
      <v-tab
        v-else
        :key="`tab-head-default-${tab.id}`"
        :class="tabTabClass(tab)"
        :disabled="tab['test-query'] && !!disabledTestQueryTooltip"
      >
        <template v-if="tab.general">
          {{ $t('common.general') }}
        </template>
        <template v-else-if="tab.patterns">
          {{ patternsLabel || $tc('common.pattern') }}
        </template>
        <template v-else-if="tab.additional">
          {{ additionalLabel || $t('common.additional') }}
        </template>
        <template v-else-if="tab['test-query']">
          <v-tooltip :disabled="!disabledTestQueryTooltip" top>
            <template #activator="{ on }">
              <span v-on="on">{{ $t('common.testQuery') }}</span>
            </template>
            <span>{{ disabledTestQueryTooltip }}</span>
          </v-tooltip>
        </template>
      </v-tab>
    </template>

    <v-tab-item
      v-for="tab in visibleTabs"
      :key="`item-${tab.id}`"
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
        :form="form"
        name="test-query"
      />

      <template-testing-test-variables
        v-else-if="tab.testing"
        :general-form="form"
        :variables-fields="variablesFields"
        :template-vars="templateVars"
        :rule-id="ruleId"
        :type="type"
        :active="isActiveTestingTab"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { isNumber } from 'lodash';
import {
  computed,
  ref,
  toRef,
  useSlots,
  watch,
  onMounted,
} from 'vue';

import { FORM_GENERAL_PATTERNS_TABS } from '@/constants';

import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import {
  useTestVariablesTabData,
} from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/test-variables/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/test-variables/partials/template-testing-test-variables-tab.vue';

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
  components: { TemplateTestingTestVariables, TemplateTestingTestVariablesTab },
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
    aiChatExpandTabKeys: {
      type: Object,
      default: null,
    },
  },
  setup(props, { emit }) {
    const slots = useSlots();

    const activeTab = ref(0);

    const generalElement = ref(null);
    const patternsElement = ref(null);
    const additionalElement = ref(null);

    const setGeneralRef = refElement => generalElement.value = refElement;
    const setPatternsRef = refElement => patternsElement.value = refElement;
    const setAdditionalRef = refElement => additionalElement.value = refElement;

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasPatternsError } = useValidationElementChildren(patternsElement);
    const { hasChildrenError: hasAdditionalError } = useValidationElementChildren(additionalElement);
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

    const hasTemplateTestingTab = computed(() => isNumber(props.type) && hasAccessToTemplateTesting.value);
    const templateTestingPending = computed(() => templateVarsPending.value || copyVarsPending.value);

    const hasGeneralSlot = computed(() => Boolean(slots.general));
    const hasPatternsSlot = computed(() => Boolean(slots.patterns));
    const hasTestQuerySlot = computed(() => Boolean(slots['test-query']));
    const hasAdditionalSlot = computed(() => Boolean(slots.additional));

    const visibleTabs = computed(() => {
      const generalTab = hasGeneralSlot.value && !props.hideGeneral && createTab(FORM_GENERAL_PATTERNS_TABS.general);
      const patternsTab = hasPatternsSlot.value && createTab(FORM_GENERAL_PATTERNS_TABS.patterns);
      const testQueryTab = hasTestQuerySlot.value && createTab(FORM_GENERAL_PATTERNS_TABS.testQuery);
      const testingTab = hasTemplateTestingTab.value && createTab(FORM_GENERAL_PATTERNS_TABS.testing);
      const additionalTab = hasAdditionalSlot.value && createTab(FORM_GENERAL_PATTERNS_TABS.additional);

      const tabs = props.reverse ? [patternsTab, generalTab, additionalTab] : [generalTab, additionalTab, patternsTab];

      tabs.push(testQueryTab, testingTab);

      return tabs.filter(Boolean);
    });

    const patternsTabIndex = computed(() => (
      visibleTabs.value.findIndex(tab => tab.id === FORM_GENERAL_PATTERNS_TABS.patterns)
    ));

    const testingTabIndex = computed(() => (
      visibleTabs.value.findIndex(tab => tab.id === FORM_GENERAL_PATTERNS_TABS.testing)
    ));

    const isActiveTestingTab = computed(() => (
      testingTabIndex.value !== -1 && activeTab.value === testingTabIndex.value
    ));

    const tabTabClass = tab => ({
      'error--text': (tab.general && hasGeneralError.value)
        || (tab.patterns && hasPatternsError.value)
        || (tab.additional && hasAdditionalError.value),
      'v-tab--tooltip': tab['test-query'],
    });

    const tabItemClass = tab => ((tab.general || tab.patterns || tab.additional) ? 'pt-4' : '');

    const expandNeededTab = computed(() => {
      if (props.aiChatExpandTabKeys) {
        return Object.fromEntries(
          Object.entries(props.aiChatExpandTabKeys).map(([fieldKey, tabId]) => [
            fieldKey,
            visibleTabs.value.findIndex(tab => tab.id === tabId),
          ]),
        );
      }

      return patternsTabIndex.value;
    });

    useAiChatExpand({
      activeTab,
      neededTab: expandNeededTab,
    });

    watch(() => props.hideGeneral, (value) => {
      if (value) {
        activeTab.value = patternsTabIndex.value;
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

      visibleTabs,

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
