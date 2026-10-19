<template>
  <c-expand-btn
    :class="expandButtonClass"
    :model-value="expanded"
    :loading="pending"
    class="alarms-expand-panel-btn"
    @update:model-value="showExpandPanel"
  />
</template>

<script>
import { mapActions } from 'pinia';

import { useAlarmStore } from '@/stores/modules/alarm';

import {
  convertAlarmDetailsQueryToRequest,
  prepareAlarmDetailsQuery,
} from '@/helpers/entities/alarm/query';

export default {
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

  emits: ['update:expanded'],

  data() {
    return {
      pending: false,
    };
  },

  computed: {
    expandButtonClass() {
      return {
        'alarms-expand-panel-btn--small': this.small,
      };
    },
  },

  methods: {
    ...mapActions(useAlarmStore, ['fetchAlarmDetails']),

    async showExpandPanel(expanded) {
      if (expanded) {
        this.pending = true;

        try {
          const query = prepareAlarmDetailsQuery(this.alarm, this.widget, this.search);

          await this.fetchAlarmDetails({
            widgetId: this.widget._id,
            id: this.alarm._id,
            query: convertAlarmDetailsQueryToRequest(query),
          });
        } catch (err) {
          console.error(err);
        } finally {
          this.pending = false;
        }
      }

      this.$emit('update:expanded', expanded);
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
    max-height: 22px;
  }
}
</style>
