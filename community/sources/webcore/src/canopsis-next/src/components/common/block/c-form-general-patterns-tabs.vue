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
      :eager="tab.general || tab.patterns"
      :disabled="tab.testing && isEmptyVariablesFields"
    >
      <slot
        v-if="tab.general"
        :set-ref="setGeneralRef"
        :copy-vars="copyVars"
        :template-vars="templateVars"
        name="general"
      />
      <slot
        v-else-if="tab.patterns"
        :set-ref="setPatternsRef"
        :template-vars="templateVars"
        :copy-vars="copyVars"
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
  onMounted,
} from 'vue';

import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import {
  useTestVariablesTabData,
} from '@/components/other/template-testing/test-variables/hooks/template-test-variables-wrapper';

import TemplateTestingTestVariables from '@/components/other/template-testing/test-variables/template-testing-test-variables.vue';
import TemplateTestingTestVariablesTab from '@/components/other/template-testing/test-variables/partials/template-testing-test-variables-tab.vue';

const TAB_GENERAL = 'general';
const TAB_PATTERNS = 'patterns';
const TAB_TEST_QUERY = 'test-query';
const TAB_TESTING = 'testing';

const createTab = id => ({
  id,
  [TAB_GENERAL]: id === TAB_GENERAL,
  [TAB_PATTERNS]: id === TAB_PATTERNS,
  [TAB_TEST_QUERY]: id === TAB_TEST_QUERY,
  [TAB_TESTING]: id === TAB_TESTING,
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
    disabledTestQueryTooltip: {
      type: String,
      required: false,
      default: '',
    },
  },
  setup(props, { emit }) {
    const slots = useSlots();

    const activeTab = ref(0);

    const generalElement = ref(null);
    const patternsElement = ref(null);

    const setGeneralRef = (refElement) => { generalElement.value = refElement; };
    const setPatternsRef = (refElement) => { patternsElement.value = refElement; };

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasPatternsError } = useValidationElementChildren(patternsElement);

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

    const visibleTabs = computed(() => {
      const tabs = [];

      if (props.reverse) {
        if (hasPatternsSlot.value) {
          tabs.push(createTab(TAB_PATTERNS));
        }
        if (hasGeneralSlot.value) {
          tabs.push(createTab(TAB_GENERAL));
        }
      } else {
        if (hasGeneralSlot.value) {
          tabs.push(createTab(TAB_GENERAL));
        }
        if (hasPatternsSlot.value) {
          tabs.push(createTab(TAB_PATTERNS));
        }
      }

      if (hasTestQuerySlot.value) {
        tabs.push(createTab(TAB_TEST_QUERY));
      }

      if (hasTemplateTestingTab.value) {
        tabs.push(createTab(TAB_TESTING));
      }

      return tabs;
    });

    const patternsTabIndex = computed(() => {
      const idx = visibleTabs.value.findIndex(tab => tab.id === TAB_PATTERNS);

      return idx >= 0 ? idx : 0;
    });

    const testingTabIndex = computed(() => visibleTabs.value.findIndex(tab => tab.id === TAB_TESTING));

    const isActiveTestingTab = computed(() => (
      testingTabIndex.value !== -1 && activeTab.value === testingTabIndex.value
    ));

    const tabTabClass = tab => ({
      'error--text': (tab.general && hasGeneralError.value) || (tab.patterns && hasPatternsError.value),
      'v-tab--tooltip': tab['test-query'],
    });

    const tabItemClass = tab => ((tab.general || tab.patterns) ? 'pt-4' : '');

    useAiChatExpand({
      activeTab,
      neededTab: patternsTabIndex,
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
