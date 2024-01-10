<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <pbehaviors-calendar-component :entity-id="config.entityId" />
    </template>
  </modal-wrapper>
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
