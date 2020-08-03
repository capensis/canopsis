<template lang="pug">
  div
    v-layout(column, align-center)
      v-flex(xs2)
        v-checkbox(
          v-model="readOnly",
          label="Read only",
          color="primary"
        )
        v-btn(@click="showPlanningModal") Show pbehavior planning
</template>

<script>
import { MODALS } from '@/constants';

export default {
  data() {
    return {
      readOnly: false,
    };
  },
  methods: {
    showPlanningModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          readOnly: this.readOnly,
        },
      });
    },
  },
};
</script>
