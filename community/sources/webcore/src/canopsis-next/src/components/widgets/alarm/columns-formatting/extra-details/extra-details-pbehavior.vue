<template>
  <div>
    <c-simple-tooltip
      :content="tooltipContent"
      top
      disable-resize
    >
      <template #activator="{ on }">
        <c-alarm-pbehavior-chip
          :color="pbehaviorInfo.color"
          :icon="pbehaviorInfo.icon_name"
          v-on="on"
        />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { useExtraDetailsPbehaviorTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    pbehaviorInfo: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsPbehaviorTooltip(props);

    return { tooltipContent };
  },
};
</script>
