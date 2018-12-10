<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.alarmsList.title') }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-card-text
      alarm-list(:widget="config.widget")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';
import queryMixin from '@/mixins/query';

import AlarmList from '@/components/other/alarm/alarms-list.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmList },
  mixins: [modalInnerMixin, queryMixin],
  created() {
    this.mergeQuery({
      id: this.config.widget._id,
      query: this.config.query,
    });
  },
};
</script>
