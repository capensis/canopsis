<template>
  <v-row>
    <v-col cols="6">
      <c-enabled-field
        v-model="oldMode"
        :label="$t('common.numberField')"
        @update:model-value="updateOldMode"
      />
    </v-col>
    <v-col cols="6">
      <c-number-field
        v-if="oldMode"
        :value="value"
        :name="name"
        required
        @update:value="updateModel"
      />
      <c-alarm-state-field
        v-else
        :value="value"
        :name="name"
        @update:value="updateModel"
      />
    </v-col>
  </v-row>
</template>

<script>
import { ALARM_STATES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  emits: ['update:value'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Number,
      required: true,
    },
    name: {
      type: String,
      default: 'state',
    },
  },
  data() {
    return {
      oldMode: this.value > ALARM_STATES.critical,
    };
  },
  methods: {
    updateOldMode(value) {
      if (!value && this.value > ALARM_STATES.critical) {
        this.updateModel(ALARM_STATES.ok);
      }
    },
  },
};
</script>
