<template>
  <div>
    <c-tooltip
      v-if="tooltipContent"
      :content="tooltipContent"
      html
      top
    >
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip
          v-bind="activatorProps"
          :color="COLORS.state.pause"
          icon="pause"
        />
      </template>
    </c-tooltip>
    <c-alarm-extra-details-chip
      v-else
      :color="COLORS.state.pause"
      icon="pause"
    />
  </div>
</template>

<script>
import { COLORS } from '@/config';

import { useExtraDetailsPauseTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    pause: {
      type: Object,
      required: true,
    },
    entity: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { tooltipContent } = useExtraDetailsPauseTooltip(props);

    return {
      COLORS,
      tooltipContent,
    };
  },
};
</script>
