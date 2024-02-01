<template lang="pug">
  c-expand-btn.alarms-expand-panel-btn(
    :class="expandButtonClass",
    :expanded="expanded",
    :loading="pending",
    @expand="showExpandPanel"
  )
</template>

<script>
import { TOURS } from '@/constants';

import { getStepClass } from '@/helpers/tour';
import { prepareAlarmDetailsQuery, convertAlarmDetailsQueryToRequest } from '@/helpers/query';

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
    small: {
      type: Boolean,
      default: false,
    },
    search: {
      type: String,
      default: '',
    },
  },
  computed: {
    expandButtonClass() {
      return {
        [getStepClass(TOURS.alarmsExpandPanel, 1)]: this.isTourEnabled,
        'alarms-expand-panel-btn--small': this.small,
      };
    },
  },
  methods: {
    async showExpandPanel() {
      if (!this.expanded) {
        this.query = prepareAlarmDetailsQuery(this.alarm, this.widget, this.search);

        await this.fetchAlarmDetails({
          widgetId: this.widget._id,
          id: this.alarm._id,
          query: convertAlarmDetailsQueryToRequest(this.query),
        });
      }

      this.$emit('input', !this.expanded);
    },
  },
};
</script>

<style lang="scss">
.alarms-expand-panel-btn {
  &--small {
    width: 22px;
    max-width: 22px;
    height: 22px;
  }
}
</style>
