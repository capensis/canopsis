<template>
  <div>
    <c-form-block-row
      :label="$t('common.state')"
      :depth="depth"
      align-center
    >
      <state-criticity-field
        v-field="value.state"
        :state-values="availableStateValues"
        mandatory
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('common.note')" :depth="depth">
      <component
        :is="textareaComponent"
        v-field="value.output"
        v-validate="'required'"
        :label="label || $t('common.note')"
        :error-messages="errors.collect(outputFieldName)"
        :name="outputFieldName"
        :variables="variables"
        autofocus
      />
    </c-form-block-row>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { computed } from 'vue';

import { ALARM_STATES } from '@/constants';

import { useInfo } from '@/hooks/store/modules/info';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

export default {
  inject: ['$validator'],
  components: { StateCriticityField },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'changeState',
    },
    variables: {
      type: Array,
      default: () => [],
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { allowChangeSeverityToInfo } = useInfo();

    const outputFieldName = computed(() => `${props.name}.output`);
    const availableStateValues = computed(() => (
      allowChangeSeverityToInfo.value ? ALARM_STATES : omit(ALARM_STATES, ['ok'])
    ));

    const textareaComponent = computed(() => (
      props.variables?.length ? 'c-payload-textarea-field' : 'v-textarea'
    ));

    return {
      outputFieldName,
      availableStateValues,
      textareaComponent,
    };
  },
};
</script>
