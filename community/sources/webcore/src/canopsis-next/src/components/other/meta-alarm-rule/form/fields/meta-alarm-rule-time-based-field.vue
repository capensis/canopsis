<template>
  <c-form-block-row :label="$t('metaAlarmRule.timeInterval')">
    <v-layout class="gap-2">
      <c-duration-field
        :duration="duration"
        :label="$t('metaAlarmRule.timeInterval')"
        required
        @input="handleInput"
      />
      <c-help-icon
        :text="$t('metaAlarmRule.timeIntervalHelpText')"
        icon="help"
        max-width="300"
        left
      />
    </v-layout>
  </c-form-block-row>
</template>

<script>
import { computed } from 'vue';

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
  setup(props, { emit }) {
    const duration = computed(() => props.value ?? durationToForm());

    const handleInput = (nextDuration) => {
      emit('input', nextDuration);
    };

    return {
      duration,
      handleInput,
    };
  },
};
</script>
