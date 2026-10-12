<template>
  <div>
    <c-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-pbehavior-chip
          v-bind="activatorProps"
          :color="pbehaviorInfo.color"
          :icon="pbehaviorInfo.icon_name"
        />
      </template>
    </c-tooltip>
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
