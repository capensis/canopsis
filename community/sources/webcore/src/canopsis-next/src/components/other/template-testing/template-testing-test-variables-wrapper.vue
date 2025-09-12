<template>
  <v-tabs fixed-tabs>
    <v-tab>{{ $t('common.general') }}</v-tab>
    <v-tooltip :disabled="!disabledTestVariablesTab" top>
      <template #activator="{ on }">
        <span class="test-variables-tab__wrapper" v-on="on">
          <v-tab :disabled="disabledTestVariablesTab">
            {{ $tc('templateTesting.testVariables') }}
          </v-tab>
        </span>
      </template>
      <span>{{ $t('templateTesting.testVariablesDisabledTooltip') }}</span>
    </v-tooltip>

    <v-tab-item eager>
      <slot :template-vars="templateVars" />
    </v-tab-item>

    <v-tab-item :disabled="disabledTestVariablesTab">
      <template-testing-test-variables
        :general-form="form"
        :variables-fields="variablesFields"
        :is-new="isNew"
        :type="type"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { computed, onMounted, toRef } from 'vue';

import TemplateTestingTestVariables from '@/components/other/template-testing/template-testing-test-variables.vue';

import { useTemplateVarsList, useTestVariablesFields } from './hooks/template-test-variables-wrapper';

export default {
  components: { TemplateTestingTestVariables },
  inheritAttrs: false,
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    type: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { items: variablesFields } = useTestVariablesFields(props, emit);
    const { templateVars, pending, fetchList } = useTemplateVarsList({
      type: toRef(props, 'type'),
    });

    const disabledTestVariablesTab = computed(() => !variablesFields.value.length);

    onMounted(fetchList);

    return {
      variablesFields,
      templateVars,
      pending,
      disabledTestVariablesTab,
    };
  },
};
</script>

<style lang="scss" scoped>
.test-variables-tab__wrapper {
  display: flex;
  flex: 1 1 auto;
  width: 100%;
  max-width: 360px;
}
</style>
