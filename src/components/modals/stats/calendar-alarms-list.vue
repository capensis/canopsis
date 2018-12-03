<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.calendarAlarmsList.title') }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-card-text
      alarm-list(:widget="config.widget")
</template>

<script>
import pick from 'lodash/pick';

import { MODALS } from '@/constants';
import AlarmList from '@/components/other/alarm/alarms-list.vue';
import modalInnerMixin from '@/mixins/modal/inner';
import queryMixin from '@/mixins/query';

export default {
  name: MODALS.calendarAlarmsList,
  components: { AlarmList },
  mixins: [modalInnerMixin, queryMixin],
  created() {
    const query = pick(this.config.query, ['tstart', 'tstop']);

    if (query.tstart || query.tstop) {
      query.interval = this.$constants.LIVE_REPORTING_INTERVALS.custom;
    }

    this.mergeQuery({
      id: this.config.widget._id,
      query,
    });
  },
};
</script>
