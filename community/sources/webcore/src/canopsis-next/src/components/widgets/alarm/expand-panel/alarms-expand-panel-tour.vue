<template lang="pug">
  v-tour(:name="tourName", :steps="tourSteps", :callbacks="tourCallbacks")
</template>

<script>
import { TOURS } from '@/constants';

import { getStepTarget } from '@/helpers/tour';

import { tourInnerMixin } from '@/mixins/tour/inner';

export default {
  mixins: [tourInnerMixin],
  computed: {
    tourName() {
      return TOURS.alarmsExpandPanel;
    },
    tourSteps() {
      const stepsContentPrefix = `tours.${this.tourName}`;

      return [
        {
          target: getStepTarget(this.tourName, 1),
          content: this.$t(`${stepsContentPrefix}.step1`),
          offset: -250,
          params: {
            placement: 'right',
          },
        },
        {
          target: getStepTarget(this.tourName, 2),
          content: this.$t(`${stepsContentPrefix}.step2`),
          offset: -250,
          params: {
            placement: 'top',
          },
        },
        {
          target: getStepTarget(this.tourName, 3),
          content: this.$t(`${stepsContentPrefix}.step3`),
          offset: -250,
          params: {
            placement: 'top',
          },
        },
      ];
    },
  },
};
</script>
