<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />

    <v-tabs v-model="activeTab" fixed-tabs>
      <v-tab
        :value="TEMPLATE_TESTING_TEST_VARIABLES_TABS.general"
        :class="{ 'text-error': hasGeneralError }"
      >
        {{ t('common.general') }}
      </v-tab>

      <template-testing-test-variables-tab
        v-if="hasAccess"
        :value="TEMPLATE_TESTING_TEST_VARIABLES_TABS.testVariables"
        :disabled="isEmptyVariablesFields"
        :error="hasTestingError"
      />
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item :value="TEMPLATE_TESTING_TEST_VARIABLES_TABS.general" class="pt-2" eager>
        <slot :template-vars="templateVars" :copy-vars="copyVars" :set-ref="setGeneralRef" />
      </v-window-item>

      <v-window-item
        v-if="hasAccess"
        :value="TEMPLATE_TESTING_TEST_VARIABLES_TABS.testVariables"
        :disabled="isEmptyVariablesFields"
      >
        <template-testing-test-variables
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
import { computed, onMounted, ref, toRef } from 'vue';

import { TEMPLATE_TESTING_TEST_VARIABLES_TABS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';
import { useCopyVarsList } from '@/hooks/vars/copy';
import { useTemplateVarsList } from '@/hooks/vars/template';

import { useTestVariablesTabData } from './hooks/template-test-variables-wrapper';
import TemplateTestingTestVariablesTab from './partials/template-testing-test-variables-tab.vue';
import TemplateTestingTestVariables from './template-testing-test-variables.vue';

export default {
  components: { TemplateTestingTestVariables, TemplateTestingTestVariablesTab },
  inheritAttrs: false,

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
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const activeTab = ref(TEMPLATE_TESTING_TEST_VARIABLES_TABS.general);
    const generalElement = ref(null);
    const testingElement = ref(null);

    const setGeneralRef = refElement => (generalElement.value = refElement);
    const setTestingRef = refElement => (testingElement.value = refElement);

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasTestingError } = useValidationElementChildren(testingElement);

    const isActiveTestingTab = computed(
      () => activeTab.value === TEMPLATE_TESTING_TEST_VARIABLES_TABS.testVariables,
    );

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
      hasAccess,
      items: variablesFields,
      isEmptyItems: isEmptyVariablesFields,
    } = useTestVariablesTabData(props, toRef(props, 'type'), emit);

    const pending = computed(() => templateVarsPending.value || copyVarsPending.value);

    onMounted(() => {
      fetchCopyVarsList();
      fetchTemplateVarsList();
    });

    return {
      t,
      activeTab,
      isActiveTestingTab,
      setGeneralRef,
      setTestingRef,
      hasGeneralError,
      hasTestingError,
      copyVars,
      templateVars,
      pending,
      hasAccess,
      variablesFields,
      isEmptyVariablesFields,
      TEMPLATE_TESTING_TEST_VARIABLES_TABS,
    };
  },
};
</script>
