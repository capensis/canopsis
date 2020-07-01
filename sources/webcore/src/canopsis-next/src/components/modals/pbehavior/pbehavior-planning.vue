<template lang="pug">
  modal-wrapper(data-test="pbehaviorListModal")
    template(slot="title")
      span {{ $t('modals.pbehaviorPlanning.title') }}
    template(slot="text")
      pbehavior-planning-calendar(:pbehaviors="pbehaviors", :readOnly="readOnly")
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorPlanning,
  components: { PbehaviorPlanningCalendar, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  computed: {
    pbehaviors() {
      return this.config.pbehaviors;
    },

    readOnly() {
      return !!this.config.readOnly;
    },
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action();
      }

      this.$modals.hide();
    },
  },
};
</script>
