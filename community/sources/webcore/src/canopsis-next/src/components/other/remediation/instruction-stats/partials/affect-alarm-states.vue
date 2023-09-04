<template lang="pug">
  div.d-flex
    c-state-count-changes-chips(:affect="alarmStates.minor", :color="minorColor")
    c-state-count-changes-chips.ml-2(:affect="alarmStates.major", :color="majorColor")
    c-state-count-changes-chips.ml-2(:affect="alarmStates.critical", :color="criticalColor")
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

import { getImpactStateColor } from '@/helpers/entities/entity/color';

export default {
  props: {
    alarmStates: {
      type: Object,
      required: true,
    },
  },
  computed: {
    minorColor() {
      return getImpactStateColor(ENTITIES_STATES.minor);
    },

    majorColor() {
      return getImpactStateColor(ENTITIES_STATES.major);
    },

    criticalColor() {
      return getImpactStateColor(ENTITIES_STATES.critical);
    },
  },
};
</script>
