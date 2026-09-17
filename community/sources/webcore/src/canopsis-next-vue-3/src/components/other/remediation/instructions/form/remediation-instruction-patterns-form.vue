<template>
  <div>
    <c-patterns-field
      v-model:value="patternsField"
      :alarm-attributes="alarmAttributes"
      :entity-attributes="entityAttributes"
      :pending="pending"
      with-alarm
      with-entity
      some-required
    />
    <c-collapse-panel :title="t('remediation.pattern.tabs.pbehaviorTypes.title')" class="mt-3">
      <remediation-patterns-pbehavior-types-form v-model:form="patternsField" />
    </c-collapse-panel>
  </div>
</template>

<script>
import { useFormModel } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import RemediationPatternsPbehaviorTypesForm from '@/components/other/remediation/patterns/form/remediation-patterns-pbehavior-types-form.vue';

export default {
  components: { RemediationPatternsPbehaviorTypesForm },
  props: {
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { fetchInstructionPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();
    const patternsField = useFormModel('patterns');

    const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
      fetchInstructionPatternFields,
      props.readonly,
    );

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      t,
      pending,
      alarmAttributes,
      entityAttributes,
      patternsField,
    };
  },
};
</script>
