<template>
  <v-tabs
    v-model="activeTab"
    slider-color="primary"
    centered
  >
    <v-tab
      v-if="reverse && hasPatternsSlot"
      key="patterns"
      :class="{ 'error--text': hasPatternsError }"
    >
      {{ patternsLabel || $tc('common.pattern') }}
    </v-tab>
    <v-tab
      key="general"
      :class="{ 'error--text': hasGeneralError }"
    >
      {{ $t('common.general') }}
    </v-tab>
    <v-tab
      v-if="!reverse && hasPatternsSlot"
      key="patterns"
      :class="{ 'error--text': hasPatternsError }"
    >
      {{ patternsLabel || $tc('common.pattern') }}
    </v-tab>
    <v-tab v-if="hasTestQuerySlot">{{ $t('common.testQuery') }}</v-tab>
    <template-testing-test-variables-tab
      v-if="hasTemplateTestingTab"
      :disabled="isEmptyVariablesFields"
    />

    <v-tab-item
      v-if="reverse && hasPatternsSlot"
      key="patternsItem"
      class="pt-4"
      eager
    >
      <slot
        :set-ref="setPatternsRef"
        :template-vars="templateVars"
        :copy-vars="copyVars"
        name="patterns"
      />
    </v-tab-item>

    <v-tab-item
      key="generalItem"
      class="pt-4"
      eager
    >
      <slot
        :set-ref="setGeneralRef"
        :copy-vars="copyVars"
        :template-vars="templateVars"
        name="general"
      />
    </v-tab-item>

    <v-tab-item
      v-if="!reverse && hasPatternsSlot"
      key="patternsItem"
      class="pt-4"
      eager
    >
      <slot
        :set-ref="setPatternsRef"
        :template-vars="templateVars"
        :copy-vars="copyVars"
        name="patterns"
      />
    </v-tab-item>

    <v-tab-item v-if="hasTestQuerySlot">
      <slot :form="form" name="test-query" />
    </v-tab-item>

    <v-tab-item
      v-if="hasTemplateTestingTab"
      :disabled="isEmptyVariablesFields"
    >
      <template-testing-test-variables
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
  },
  setup(props, { emit }) {
    const slots = useSlots();

    const activeTab = ref(0);

    const generalElement = ref(null);
    const patternsElement = ref(null);

    const setGeneralRef = refElement => generalElement.value = refElement;
    const setPatternsRef = refElement => patternsElement.value = refElement;

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

    const hasTestQuerySlot = computed(() => Boolean(slots.testQuery));

    const GENERAL_PATTERNS_FORM_TABS = computed(() => {
      const result = props.reverse ? { patterns: 0, general: 1 } : { general: 0, patterns: 1 };

      if (hasTestQuerySlot.value) {
        result.testQuery = 2;
      }

      result.testing = 2 + Number(hasTestQuerySlot.value);

      return result;
    });

    const isActiveTestingTab = computed(() => activeTab.value === GENERAL_PATTERNS_FORM_TABS.value.testing);

    const hasPatternsSlot = computed(() => Boolean(slots.patterns));

    useAiChatExpand({
      activeTab,
      neededTab: GENERAL_PATTERNS_FORM_TABS.value.patterns,
    });

    onMounted(() => {
      if (hasTemplateTestingTab.value) {
        fetchCopyVarsList();
        fetchTemplateVarsList();
      }
    });

    return {
      activeTab,

      hasPatternsSlot,
      hasTestQuerySlot,
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
    };
  },
};
</script>
