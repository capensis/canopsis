<template>
  <c-expand-btn
    :class="expandButtonClass"
    :expanded="expanded"
    :loading="pending"
    class="alarms-expand-panel-btn"
    @expand="showExpandPanel"
  />
</template>

<script>
import { prepareAlarmDetailsQuery, convertAlarmDetailsQueryToRequest } from '@/helpers/entities/alarm/query';

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
    height: 22px;
    max-width: 22px;
    max-height: 22px;
  }
}
</style>
