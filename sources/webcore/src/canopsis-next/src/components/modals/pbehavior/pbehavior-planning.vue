<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.pbehaviorPlanning.title') }}
      template(slot="text")
        pbehavior-planning-calendar(
          :pbehaviorsById.sync="form.pbehaviorsById",
          :addedPbehaviorsById.sync="form.addedPbehaviorsById",
          :changedPbehaviorsById.sync="form.changedPbehaviorsById",
          :removedPbehaviorsById.sync="form.removedPbehaviorsById",
          :readOnly="readOnly",
          :filter="filter"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { keyBy, omit } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import { pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorPlanning,
  components: { PbehaviorPlanningCalendar, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin(), entitiesPbehaviorMixin],
  data() {
    return {
      form: {
        pbehaviorsById: keyBy(this.modal.config.pbehaviors, '_id'),
        addedPbehaviorsById: {},
        changedPbehaviorsById: {},
        removedPbehaviorsById: {},
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
      await this.createPbehaviors(Object.values(this.form.addedPbehaviorsById)
        .map(pbehavior => pbehaviorToRequest(omit(pbehavior, ['_id']))));
      await this.updatePbehaviors(Object.values(this.form.changedPbehaviorsById).map(pbehaviorToRequest));
      await this.removePbehaviors(Object.values(this.form.removedPbehaviorsById));

      if (this.config.afterSubmit) {
        await this.config.afterSubmit();
      }

      this.$modals.hide();
    },
  },
};
</script>
