<template lang="pug">
  modal-wrapper
    template(slot="fullTitle")
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.alarmsList.title') }}
          v-btn(icon, dark, @click.native="$modals.hide")
            v-icon close
    template(slot="text")
      alarms-list-widget(:widget="config.widget")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import AlarmsListWidget from '@/components/other/alarm/alarms-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmsListWidget, ModalWrapper },
  mixins: [modalInnerMixin],
};
</script>
