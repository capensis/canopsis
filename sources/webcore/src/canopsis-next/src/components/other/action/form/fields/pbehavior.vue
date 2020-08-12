<template lang="pug">
  v-btn.primary(@click="showPbehaviorPlanningModal") {{ $t('modals.pbehaviorPlanning.title') }}
</template>

<script>

import { MODALS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  methods: {
    showPbehaviorPlanningModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          pbehavior: this.value,
          action: pbehavior => this.updateField(pbehavior),
        },
      });
    },
  },
};
</script>
