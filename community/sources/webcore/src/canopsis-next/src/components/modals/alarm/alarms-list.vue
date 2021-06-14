<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.alarmsList.title') }}
    template(slot="text")
      alarms-list-widget(:widget="config.widget")
</template>

<script>
import { MODALS } from '@/constants';

import AlarmsListWidget from '@/components/widgets/alarm/alarms-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmsListWidget, ModalWrapper },
};
</script>
