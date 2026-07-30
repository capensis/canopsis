<template>
  <v-row align="stretch" class="pattern-group-field my-0" density="compact">
    <v-col cols="auto" class="pattern-group-field__operator">
      <pattern-operator-information>{{ $t('common.and') }}</pattern-operator-information>
    </v-col>
    <v-col class="pattern-group-field__rules">
      <pattern-rules-field
        :rules="group.rules"
        :attributes="attributes"
        :disabled="disabled"
        :readonly="readonly"
        @update:rules="updateRules"
      />
    </v-col>
  </v-row>
</template>

<script>
import { formMixin } from '@/mixins/form';

import PatternRulesField from './pattern-rules-field.vue';
import PatternOperatorInformation from './pattern-operator-information.vue';

export default {
  emits: ['update:group', 'remove'],
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

<style lang="scss" scoped>
.pattern-group-field {
  width: 100%;

  &__operator {
    display: flex;
    align-items: center;
  }

  &__rules {
    min-width: 0;
  }
}
</style>
