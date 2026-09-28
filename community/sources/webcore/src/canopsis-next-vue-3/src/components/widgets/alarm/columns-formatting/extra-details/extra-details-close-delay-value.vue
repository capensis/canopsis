<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip
          v-bind="activatorProps"
          :icon="icon"
          theme-color="cancel_delay"
        />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { useExtraDetailsCloseDelayValueTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    closeDelayValue: {
      type: Number,
      required: true,
    },
  },

  setup(props) {
    const { tooltipContent } = useExtraDetailsCloseDelayValueTooltip(props);

    const icon = '$vuetify.icons.remove_clock';
    return {
      icon,
      tooltipContent,
    };
  },
};
</script>
