<template>
  <div>
    <c-patterns-field
      v-model="formModel"
      :alarm-attributes="alarmAttributes"
      :entity-attributes="entityAttributes"
      :pending="pending"
      with-alarm
      with-entity
      some-required
    />
    <c-collapse-panel :title="$t('remediation.pattern.tabs.pbehaviorTypes.title')" class="mt-3">
      <remediation-patterns-pbehavior-types-form v-model:form="formModel" />
    </c-collapse-panel>
  </div>
</template>

<script>
import { computed } from 'vue';

import RemediationPatternsPbehaviorTypesForm from '@/components/other/remediation/patterns/form/remediation-patterns-pbehavior-types-form.vue';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';
import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  components: { RemediationPatternsPbehaviorTypesForm },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:form'],
  setup(props, { emit }) {
    const { fetchInstructionPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchInstructionPatternFields,
      props.readonly,
    );
    const formModel = computed({
      get: () => props.form,
      set: value => emit('update:form', value),
    });

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      pending,
      alarmAttributes,
      entityAttributes,
      formModel,
    };
  },
};
</script>
