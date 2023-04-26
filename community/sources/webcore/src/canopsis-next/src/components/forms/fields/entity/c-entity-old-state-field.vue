<template lang="pug">
  v-layout(row)
    v-flex(xs6)
      c-enabled-field(
        v-model="oldMode",
        :label="$t('common.numberField')",
        @change="changeOldMode"
      )
    v-flex(xs6)
      v-text-field(v-if="oldMode", v-field.number="value", type="number")
      c-entity-state-field(v-else, v-field="value")
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

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
  },
  data() {
    return {
      oldMode: this.value > ENTITIES_STATES.critical,
    };
  },
  methods: {
    changeOldMode(value) {
      if (!value && this.value > ENTITIES_STATES.critical) {
        this.updateModel(ENTITIES_STATES.ok);
      }
    },
  },
};
</script>
