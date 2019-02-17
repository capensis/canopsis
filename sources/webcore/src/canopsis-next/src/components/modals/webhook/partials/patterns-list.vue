<template lang="pug">
  div
    v-layout(
    v-for="(pattern, index) in patterns",
    :key="`${getPatternString(pattern)}${index}`",
    row,
    wrap,
    align-center
    )
      v-flex(xs11)
        v-textarea(
        :value="getPatternString(pattern)",
        rows="7",
        no-resize,
        readonly,
        disabled
        )
      v-flex.text-xs-center(xs1)
        div
          v-btn(icon, @click="showEditPatternModal(index)")
            v-icon edit
        div
          v-btn(color="error", icon, @click="showRemovePatternModal(index)")
            v-icon delete
    v-btn(color="primary", @click="showCreatePatternModal") Add pattern
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    getPatternString() {
      return (pattern) => {
        if (pattern) {
          return JSON.stringify(pattern, null, 4);
        }

        return '{}';
      };
    },
  },
  methods: {
    updatePatterns(patterns) {
      return this.$emit('input', patterns);
    },

    showCreatePatternModal() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          action: pattern => this.updatePatterns(this.patterns.concat(pattern)),
        },
      });
    },

    showEditPatternModal(index) {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          action: (pattern) => {
            const newPatterns = this.patterns.map((p, i) => {
              if (i === index) {
                return pattern;
              }

              return p;
            });

            this.updatePatterns(newPatterns);
          },
        },
      });
    },

    showRemovePatternModal(index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.updatePatterns(this.patterns.filter((p, i) => i !== index)),
        },
      });
    },
  },
};
</script>
