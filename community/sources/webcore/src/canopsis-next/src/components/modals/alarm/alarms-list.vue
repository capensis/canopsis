<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ title }}
    template(#text="")
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
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.alarmsList.title');
    },
  },
  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },
};
</script>
