<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip
          v-bind="activatorProps"
          icon="center_focus_weak"
          theme-color="parent"
        />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { useExtraDetailsParentsTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    total: {
      type: Number,
      default: 0,
    },
  },

  setup(props) {
    const { tooltipContent } = useExtraDetailsParentsTooltip(props);

    return {
      tooltipContent,
    };
  },
};
</script>

<style lang="scss" scoped>
.rule-name:nth-of-type(odd) {
  color: #b5b5b5;
}
</style>
