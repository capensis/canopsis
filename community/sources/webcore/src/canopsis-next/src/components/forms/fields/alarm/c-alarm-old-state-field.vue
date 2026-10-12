<template>
  <v-layout>
    <v-flex xs6>
      <c-enabled-field
        v-model="oldMode"
        :label="$t('common.numberField')"
        @input="updateOldMode"
      />
    </v-flex>
    <v-flex xs6>
      <c-number-field
        v-if="oldMode"
        v-field="value"
        :name="name"
        required
      />
      <c-alarm-state-field
        v-else
        v-field="value"
        :name="name"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { ALARM_STATES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
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
