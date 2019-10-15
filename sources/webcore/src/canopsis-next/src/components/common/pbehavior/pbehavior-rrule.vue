<template lang="pug">
  v-card
    v-card-text
      v-layout(row)
        v-flex(xs2)
          strong {{ $t('rRule.stringLabel') }}
        v-flex(xs10)
          p.rrule-paragraph {{ rRuleString }}
      v-layout(row)
        v-flex(xs2)
          strong {{ $t('rRule.textLabel') }}
        v-flex(xs10)
          p.rrule-paragraph {{ rRuleText }}
</template>

<script>
import { rrulestr } from 'rrule';

export default {
  props: {
    rrule: {
      type: String,
      required: true,
    },
  },
  computed: {
    rRule() {
      return rrulestr(this.rrule);
    },

    rRuleString() {
      return this.rRule.toString();
    },

    rRuleText() {
      return this.rRule.toText();
    },
  },
};
</script>

<style scoped>
  .rrule-paragraph {
    word-wrap: break-word;
  }
</style>
