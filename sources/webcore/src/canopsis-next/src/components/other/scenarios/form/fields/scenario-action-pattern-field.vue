<template lang="pug">
  div
    c-patterns-field(v-field="patterns", :disabled="disabled", :name="name", alarm, entity)
</template>

<script>
import { isEmpty } from 'lodash';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'patterns',
    },
  },
  mounted() {
    this.$validator.attach({
      name: this.name,
      rules: 'required:true',
      getter: () => !isEmpty(this.patterns.alarm_patterns)
        || !isEmpty(this.patterns.entity_patterns),
      context: () => this,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
};
</script>
