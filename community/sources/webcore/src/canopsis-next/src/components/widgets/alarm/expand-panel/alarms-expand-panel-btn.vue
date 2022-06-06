<template lang="pug">
  c-expand-btn(
    :class="expandButtonClass",
    :expanded="expanded",
    :disabled="pending",
    :loading="pending",
    @expand="showExpandPanel"
  )
</template>

<script>
import { TOURS } from '@/constants';

import { getStepClass } from '@/helpers/tour';
import { prepareAlarmDetailsQuery } from '@/helpers/query';

import { widgetExpandPanelAlarmDetails } from '@/mixins/widget/expand-panel/alarm/details';

export default {
  inject: ['$system'],
  mixins: [widgetExpandPanelAlarmDetails],
  model: {
    prop: 'expanded',
    event: 'input',
  },
  props: {
    expanded: {
      type: Boolean,
      default: false,
    },
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    expandButtonClass() {
      if (this.isTourEnabled) {
        return getStepClass(TOURS.alarmsExpandPanel, 1);
      }

      return '';
    },
  },
  methods: {
    async showExpandPanel() {
      if (!this.expanded) {
        const query = prepareAlarmDetailsQuery(this.alarm, this.widget);

        this.updateQuery({
          id: this.queryId,

          query,
        });

        await this.fetchAlarmItemDetails({
          id: this.queryId,

          query,
        });
      }

      this.$emit('input', !this.expanded);
    },
  },
};
</script>

<style lang="scss" scoped>
.not-filtered {
  opacity: .4;
  transition: opacity .3s linear;

  &:hover {
    opacity: 1;
  }
}
</style>
