<template>
  <v-row class="junit-state-setting-threshold-state-field align-center">
    <v-col cols="4" class="d-flex align-center">
      {{ label }}
    </v-col>
    <v-col cols="4" class="d-flex align-center">
      <c-number-field
        :value="modelValue"
        :name="name"
        :min="0"
        :max="max"
        class="junit-state-setting-threshold-state-field__input"
        @update:value="$emit('update:modelValue', $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { JUNIT_STATE_SETTING_THRESHOLDS_TYPES } from '@/constants';

export default {
  props: {
    modelValue: {
      type: [Number, String],
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    type: {
      type: Number,
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    name: {
      type: String,
      default: 'state',
    },
  },

  emits: ['update:modelValue'],

  setup(props) {
    const max = computed(() =>
      props.type === JUNIT_STATE_SETTING_THRESHOLDS_TYPES.percent ? 100 : undefined,
    );

    return { max };
  },
};
</script>

<style lang="scss" scoped>
.junit-state-setting-threshold-state-field {
  min-height: 48px;

  &__input {
    width: 100%;
  }
}
</style>
