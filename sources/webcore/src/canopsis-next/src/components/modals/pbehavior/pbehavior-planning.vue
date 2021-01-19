<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(fill-height, close)
      template(slot="title")
        span {{ $t('modals.pbehaviorPlanning.title') }}
      template(slot="text")
        pbehavior-planning-calendar(
          :pbehaviors-by-id.sync="form.pbehaviorsById",
          :added-pbehaviors-by-id.sync="form.addedPbehaviorsById",
          :changed-pbehaviors-by-id.sync="form.changedPbehaviorsById",
          :removed-pbehaviors-by-id.sync="form.removedPbehaviorsById",
          :read-only="readOnly",
          :filter="filter"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", :loading="submitting", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { keyBy, omit } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import { pbehaviorToDuplicateForm, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';

import PbehaviorPlanningCalendar from '@/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorPlanning,
  components: { PbehaviorPlanningCalendar, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesPbehaviorMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { pbehaviors = [], pbehaviorsToAdd = [] } = this.modal.config;

    return {
      form: {
        pbehaviorsById: keyBy(pbehaviors, '_id'),
        addedPbehaviorsById: keyBy(pbehaviorsToAdd.map(pbehaviorToDuplicateForm), '_id'),
        changedPbehaviorsById: {},
        removedPbehaviorsById: {},
      },
    };
  },
  computed: {
    readOnly() {
      return !!this.config.readOnly;
    },

    filter() {
      return this.config.filter;
    },
  },
  methods: {
    async submit() {
      try {
        await this.createPbehaviors(Object.values(this.form.addedPbehaviorsById)
          .map(pbehavior => pbehaviorToRequest(omit(pbehavior, ['_id']))));
        await this.updatePbehaviors(Object.values(this.form.changedPbehaviorsById).map(pbehaviorToRequest));
        await this.removePbehaviors(Object.values(this.form.removedPbehaviorsById));

        if (this.config.afterSubmit) {
          await this.config.afterSubmit();
        }

        this.$modals.hide();
      } catch (err) {
        const message = Object.values(err).join(' ');

        this.$popups.error({ text: message || this.$t('errors.default') });
      }
    },
  },
};
</script>
