<template>
  <div>
    <c-form-block-row
      :label="t('common.state')"
      :depth="depth"
      align-center
    >
      <state-criticity-field
        :model-value="selectedValue.state"
        :state-values="availableStateValues"
        mandatory
        @update:model-value="updateField('state', $event)"
      />
    </c-form-block-row>

    <c-form-block-row :label="t('common.note')" :depth="depth">
      <component
        :is="textareaComponent"
        :model-value="selectedValue.output"
        :label="label || t('common.note')"
        :name="outputFieldName"
        :variables="variables"
        required
        autofocus
        @update:model-value="updateField('output', $event)"
        @update:value="updateField('output', $event)"
      />
    </c-form-block-row>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { computed } from 'vue';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';
import { ALARM_STATES } from '@/constants';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

export default {
  components: { StateCriticityField },

  props: {
    modelValue: {
      type: Object,
      default: undefined,
    },
    value: {
      type: Object,
      default: () => ({}),
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

  emits: ['update:modelValue', 'update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { allowChangeSeverityToInfo } = useInfo();
    const { updateField } = useModelField(props, emit);

    const selectedValue = computed(() => props.modelValue ?? props.value ?? {});
    const outputFieldName = computed(() => `${props.name}.output`);
    const availableStateValues = computed(() => (
      allowChangeSeverityToInfo.value ? ALARM_STATES : omit(ALARM_STATES, ['ok'])
    ));
    const textareaComponent = computed(() => (
      props.variables?.length ? 'c-payload-textarea-field' : 'c-textarea'
    ));

    return {
      t,
      selectedValue,
      outputFieldName,
      availableStateValues,
      textareaComponent,
      updateField,
    };
  },
};
</script>
