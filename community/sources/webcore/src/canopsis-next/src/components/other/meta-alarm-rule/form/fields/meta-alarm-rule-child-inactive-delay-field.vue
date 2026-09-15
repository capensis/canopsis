<template>
  <c-form-block-row :label="$t('metaAlarmRule.childInactiveDelay')">
    <v-layout class="gap-3" column>
      <v-layout class="gap-2">
        <c-enabled-field
          :value="!!value"
          :label="$t('metaAlarmRule.childInactiveDelay')"
          @input="handleToggle"
        />
        <c-help-icon
          :text="$t('metaAlarmRule.childInactiveDelayHelpText')"
          icon="help"
          max-width="300"
          top
        />
      </v-layout>
      <c-duration-field
        v-if="value"
        :duration="value"
        clearable
        @input="handleDurationInput"
      />
    </v-layout>
  </c-form-block-row>
</template>

<script>
import { durationToForm } from '@/helpers/date/duration';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: undefined,
    },
  },
  setup(_, { emit }) {
    const handleToggle = (enabled) => {
      emit('input', enabled ? durationToForm() : undefined);
    };

    const handleDurationInput = (duration) => {
      emit('input', duration);
    };

    return {
      handleToggle,
      handleDurationInput,
    };
  },
};
</script>
