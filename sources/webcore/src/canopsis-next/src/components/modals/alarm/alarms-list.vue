<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.alarmsList.title') }}
        v-btn(icon, dark, @click.native="$modals.hide")
          v-icon close
    v-card-text
      alarms-list-widget(:widget="config.widget")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import queryMixin from '@/mixins/query';

import AlarmsListWidget from '@/components/other/alarm/alarms-list.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmsListWidget },
  mixins: [modalInnerMixin, queryMixin],
  created() {
    this.mergeQuery({
      id: this.config.widget._id,
      query: this.config.query,
    });
  },
};
</script>
