<template lang="pug">
  v-tour(:name="tourName", :steps="tourSteps", :callbacks="tourCallbacks")
</template>

<script>
import { TOURS } from '@/constants';

import tourInnerMixin from '@/mixins/tour/inner';

export default {
  mixins: [tourInnerMixin],
  computed: {
    tourName() {
      return TOURS.alarmsExpand;
    },
    tourSteps() {
      return [
        {
          target: '.v-tour-step-1',
          content: 'Step 1',
          offset: -250,
          params: {
            placement: 'right',
          },
        },
        {
          target: '.v-tour-step-2',
          content: 'Step 2',
          offset: -250,
          params: {
            placement: 'top',
          },
        },
        {
          target: '.v-tour-step-3',
          content: 'Step 3',
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
