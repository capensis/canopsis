<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      pbehaviors-calendar-component(:entity-id="config.entityId")
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import PbehaviorsCalendarComponent from '@/components/other/pbehavior/calendar/pbehaviors-calendar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorsCalendar,
  components: { PbehaviorsCalendarComponent, ModalWrapper },
  mixins: [modalInnerMixin],
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.pbehaviorsCalendar.title');
    },
  },
};
</script>
