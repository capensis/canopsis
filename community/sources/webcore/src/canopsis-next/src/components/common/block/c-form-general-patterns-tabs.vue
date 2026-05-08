<template>
  <v-tabs
    v-model="activeTab"
    slider-color="primary"
    centered
  >
    <v-tab
      key="patterns"
      :class="{ 'error--text': hasPatternsError }"
    >
      {{ $tc('common.pattern') }}
    </v-tab>
    <v-tab
      key="general"
      :class="{ 'error--text': hasGeneralError }"
    >
      {{ $t('common.general') }}
    </v-tab>
    <template-testing-test-variables-tab
      v-if="hasTemplateTestingTab"
      :disabled="isEmptyVariablesFields"
    />

    <v-tab-item
      key="patternsItem"
      class="pt-4"
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
    >
      <slot
        :set-ref="setGeneralRef"
        :copy-vars="copyVars"
        :template-vars="templateVars"
        name="general"
      />
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
  watch,
  onMounted,
} from 'vue';

import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';
import { useAiChatExpand } from '@/hooks/ai/ai-chat-form';

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
  },
  setup(props, { emit }) {
    const slots = useSlots();

    const GENERAL_PATTERNS_FORM_TABS = computed(() => ({
      patterns: 0,
      general: 1,
      testing: 2,
    }));

    const activeTab = ref(0);

    const isActiveTestingTab = computed(() => activeTab.value === GENERAL_PATTERNS_FORM_TABS.value.testing);

    const hasPatternsSlot = computed(() => Boolean(slots.patterns));

    useAiChatExpand({
      activeTab,
      neededTab: GENERAL_PATTERNS_FORM_TABS.value.patterns,
    });

    const generalElement = ref(null);
    const patternsElement = ref(null);

    const hasGeneralError = ref(false);
    const hasPatternsError = ref(false);

    const setGeneralRef = refElement => generalElement.value = refElement;
    const setPatternsRef = refElement => patternsElement.value = refElement;

    watch(() => generalElement.value?.hasAnyError, (value) => {
      hasGeneralError.value = value ?? false;
    });

    watch(() => patternsElement.value?.hasAnyError, (value) => {
      hasPatternsError.value = value ?? false;
    });

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

    onMounted(() => {
      if (hasTemplateTestingTab.value) {
        fetchCopyVarsList();
        fetchTemplateVarsList();
      }
    });

    return {
      activeTab,

      hasPatternsSlot,

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
