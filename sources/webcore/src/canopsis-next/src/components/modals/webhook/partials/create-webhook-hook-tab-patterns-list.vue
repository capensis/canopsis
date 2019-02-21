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
    v-btn(color="primary", @click="showCreatePatternModal") {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formArrayMixin from '@/mixins/form/array';

export default {
  mixins: [modalMixin, formArrayMixin],
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
    showCreatePatternModal() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          action: pattern => this.addItemIntoArray(pattern),
        },
      });
    },

    showEditPatternModal(index) {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          pattern: this.patterns[index],
          action: pattern => this.updateItemInArray(index, pattern),
        },
      });
    },

    showRemovePatternModal(index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>
