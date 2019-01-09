<template lang="pug">
  div
    v-layout(justify-end)
      v-btn(@click="showAddRuleFieldModal", small) Add a field
    div(v-for="(rule, ruleKey) in pattern", :key="ruleKey")
      v-card.my-1(flat, dark, tile)
        v-card-text
          v-layout
            v-flex
              p Field : {{ ruleKey }}
            v-flex(v-if="isSimpleRule(rule)")
              p Value : {{ rule }}
            v-flex(v-else)
              v-layout(column)
                v-flex(v-for="(field, fieldKey) in rule", :key="fieldKey")
                  p {{ fieldKey }} {{ field }}
            v-flex
              v-layout(justify-end)
                v-btn(@click="editRule(ruleKey)", icon, small)
                  v-icon edit
                v-btn(@click="deleteRule(ruleKey)", icon, small)
                  v-icon(color="error") delete
</template>

<script>
import omit from 'lodash/omit';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';
import modalMixin from '@/mixins/modal';

export default {
  mixins: [formMixin, modalMixin],
  model: {
    prop: 'pattern',
    event: 'input',
  },
  props: {
    pattern: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      required: true,
    },
  },
  methods: {
    isSimpleRule(rule) {
      if (typeof rule === 'string') {
        return true;
      }

      return false;
    },

    deleteRule(rule) {
      this.removeField(rule);
    },

    editRule(rule) {
      this.showModal({
        name: MODALS.addEventFilterRuleToPattern,
        config: {
          ruleKey: rule,
          ruleValue: this.pattern[rule],
          isSimpleRule: this.isSimpleRule(this.pattern[rule]),
          operators: this.operators,
          action: (newRule) => {
            const newPattern = omit(this.pattern, rule);
            newPattern[newRule.field] = newRule.value;
            this.$emit('input', newPattern);
          },
        },
      });
    },

    showAddRuleFieldModal() {
      this.showModal({
        name: MODALS.addEventFilterRuleToPattern,
        config: {
          operators: this.operators,
          action: newRule => this.updateField(newRule.field, newRule.value),
        },
      });
    },
  },
};
</script>
