<template>
  <c-layout>
    <div class="pa-2 mr-2">
      <pattern-operator-information>{{ $t('common.and') }}</pattern-operator-information>
    </div>
    <div :style="rulesColumnStyle">
      <c-layout column>
        <pattern-rules-field
          :rules="group.rules"
          :attributes="attributes"
          :disabled="disabled"
          :readonly="readonly"
          @update:rules="updateRules"
        />
      </c-layout>
    </div>
  </c-layout>
</template>

<script>
import { gridFlexStyle } from '@/helpers/layout';

import { formMixin } from '@/mixins/form';

import PatternRulesField from './pattern-rules-field.vue';
import PatternOperatorInformation from './pattern-operator-information.vue';

export default {
  emits: ['update:group'],
  components: { PatternOperatorInformation, PatternRulesField },
  mixins: [formMixin],
  props: {
    group: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rulesColumnStyle() {
      return gridFlexStyle(11);
    },
  },
  methods: {
    updateRules(rules) {
      if (rules.length) {
        this.updateField('rules', rules);
      } else {
        this.$emit('remove');
      }
    },
  },
};
</script>
