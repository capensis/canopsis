<template lang="pug">
  model-wrapper
    template(slot="fullTitle")
      span.headline {{ $t('modals.alarmsList.title') }}
      v-btn(icon, dark, @click.native="$modals.hide")
        v-icon close
    template(slot="text")
      alarms-list-widget(:widget="config.widget")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import queryMixin from '@/mixins/query';

import AlarmsListWidget from '@/components/other/alarm/alarms-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmsListWidget, ModalWrapper },
  mixins: [modalInnerMixin, queryMixin],
  mounted() {
    this.mergeQuery({
      id: this.config.widget._id,
      query: this.config.query,
    });
  },
};
</script>
