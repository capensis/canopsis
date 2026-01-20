<template>
  <div>
    <c-patterns-field
      v-field="form"
      :alarm-attributes="alarmAttributes"
      :entity-attributes="entityAttributes"
      :pending="pending"
      with-alarm
      with-entity
      both-counters
      some-required
    />
    <c-collapse-panel
      :title="$t('remediation.pattern.tabs.pbehaviorTypes.title')"
      class="mt-3"
    >
      <remediation-patterns-pbehavior-types-form v-field="form" />
    </c-collapse-panel>
  </div>
</template>

<script>
import { useValidationHeader } from '@/hooks/validator/validation-header';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

import RemediationPatternsPbehaviorTypesForm from '@/components/other/remediation/patterns/form/remediation-patterns-pbehavior-types-form.vue';

export default {
  components: { RemediationPatternsPbehaviorTypesForm },
  model: {
    prop: 'form',
    event: 'input',
  },
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
  setup() {
    const { fetchInstructionPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchInstructionPatternFields);

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      pending,
      alarmAttributes,
      entityAttributes,
    };
  },
};
</script>
