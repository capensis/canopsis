<template>
  <div class="d-flex">
    <c-state-count-changes-chips
      :affect="alarmStates.minor"
      :color="minorColor"
    />
    <c-state-count-changes-chips
      class="ml-2"
      :affect="alarmStates.major"
      :color="majorColor"
    />
    <c-state-count-changes-chips
      class="ml-2"
      :affect="alarmStates.critical"
      :color="criticalColor"
    />
  </div>
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

import { getEntityStateColor } from '@/helpers/entities/entity/color';

export default {
  props: {
    alarmStates: {
      type: Object,
      required: true,
    },
  },
  computed: {
    minorColor() {
      return getEntityStateColor(ENTITIES_STATES.minor);
    },

    majorColor() {
      return getEntityStateColor(ENTITIES_STATES.major);
    },

    criticalColor() {
      return getEntityStateColor(ENTITIES_STATES.critical);
    },
  },
};
</script>
