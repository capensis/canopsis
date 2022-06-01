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
import { prepareAlarmDetailsQuery, generateAlarmDetailsQueryId } from '@/helpers/query';

import { queryMixin } from '@/mixins/query';
import { entitiesAlarmDetailsMixin } from '@/mixins/entities/alarm/details';
import { getStepClass } from '@/helpers/tour';
import { TOURS } from '@/constants';

export default {
  inject: ['$system'],
  mixins: [queryMixin, entitiesAlarmDetailsMixin],
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
  },
  data() {
    return {
      pending: false,
    };
  },
  computed: {
    expandButtonClass() {
      if (this.isTourEnabled) { // TODO: move this logic to mixin
        return getStepClass(TOURS.alarmsExpandPanel, 1);
      }

      return '';
    },
  },
  methods: {
    async showExpandPanel() {
      if (!this.expanded) {
        this.pending = true;
        const query = prepareAlarmDetailsQuery(this.alarm, this.widget);

        this.updateQuery({
          id: generateAlarmDetailsQueryId(this.alarm, this.widget),

          query,
        });

        await this.fetchAlarmItemDetails({
          data: [query],
        });

        this.pending = false;
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
