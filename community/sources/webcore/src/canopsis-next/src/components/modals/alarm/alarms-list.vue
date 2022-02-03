<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.alarmsList.title') }}
    template(slot="text")
      alarms-list-widget(:widget="config.widget", local-widget)
</template>

<script>
import { MODALS } from '@/constants';

import Observer from '@/services/observer';

import { modalInnerMixin } from '@/mixins/modal/inner';

import AlarmsListWidget from '@/components/widgets/alarm/alarms-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.alarmsList,
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
  components: { AlarmsListWidget, ModalWrapper },
  mixins: [modalInnerMixin],
  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },
};
</script>
