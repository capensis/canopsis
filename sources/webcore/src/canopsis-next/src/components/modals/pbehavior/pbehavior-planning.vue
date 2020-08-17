<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.pbehaviorPlanning.title') }}
      template(slot="text")
        pbehavior-planning-calendar(v-model="form", :readOnly="readOnly", :filter="filter")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorPlanning,
  components: { PbehaviorPlanningCalendar, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin(), entitiesPbehaviorMixin],
  data() {
    return {
      form: {
        list: this.modal.config.pbehaviors || [],
        added: [],
        changed: [],
        removed: [],
      },
    };
  },
  computed: {
    pbehaviors() {
      return this.config.pbehaviors;
    },

    readOnly() {
      return !!this.config.readOnly;
    },

    filter() {
      return this.config.filter;
    },
  },
  methods: {
    async submit() {
      await this.createPbehaviors(this.form.added);
      await this.removePbehaviors(this.form.removed);
      await this.updatePbehaviors(this.form.changed);

      if (this.config.afterSubmit) {
        await this.config.afterSubmit();
      }

      this.$modals.hide();
    },
  },
};
</script>
