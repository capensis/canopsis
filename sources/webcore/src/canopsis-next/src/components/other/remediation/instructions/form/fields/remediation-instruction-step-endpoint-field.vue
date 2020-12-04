<template lang="pug">
  v-layout.mt-2(column)
    v-layout.py-1
      v-flex.mt-3(xs1)
        draggable-step-number(
          :color="draggableStepNumberColor",
          disabled
        ) {{ $t('remediationInstructions.endpointAvatar') }}
      v-flex(xs11)
        v-layout(row)
          v-flex.px-1(xs11)
            v-text-field(
              v-field="value",
              v-validate="'required'",
              :label="$t('remediationInstructions.endpoint')",
              :name="name",
              :error-messages="errorMessages",
              box
            )
              v-tooltip(slot="append", left)
                v-icon(slot="activator") help
                span {{ $t('remediationInstructions.tooltips.endpoint') }}
          v-flex(xs1)
</template>

<script>
import uid from '@/helpers/uid';

import DraggableStepNumber from '../../partials/draggable-step-number.vue';

export default {
  inject: ['$validator'],
  components: { DraggableStepNumber },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  computed: {
    fieldSuffix() {
      return uid();
    },
    name() {
      return `endpoint${this.fieldSuffix}`;
    },
    errorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.fieldSuffix, ''));
    },
    draggableStepNumberColor() {
      return this.errors.has(this.name) ? 'error' : 'primary';
    },
  },
};
</script>
