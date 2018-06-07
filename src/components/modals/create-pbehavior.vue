<template lang="pug">
  v-form(@submit.prevent="submit", slot-scope="slotProps")
    v-card
      v-card-title
        span.headline {{ $t('modals.createPbehavior.title') }}
      v-card-text
        v-layout(row)
          v-text-field(
          :label="$t('modals.createPbehavior.fields.name')",
          :error-messages="errors.collect('name')",
          v-model="form.name",
          v-validate="'required'",
          data-vv-name="name"
          )
        v-layout(row)
          date-time-picker(
          :label="$t('modals.createPbehavior.fields.start')",
          v-model="form.tstart",
          name="tstart",
          rules="required",
          )
        v-layout(row)
          date-time-picker(
          :label="$t('modals.createPbehavior.fields.stop')",
          v-model="form.tstop",
          name="tstop",
          rules="required"
          )
        r-rule-form(@input="changeRRule")
        v-layout(row)
          v-select(
          label="Reason",
          v-model="form.reason",
          :items="selectItems.reasons"
          )
        v-layout(row)
          v-select(
          label="Type",
          v-model="form.type_",
          :items="selectItems.types"
          )
        v-layout(row)
          v-alert(:value="this.serverError", type="error")
            span {{ this.serverError }}
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import DateTimePicker from '@/components/forms/date-time-picker.vue';
import RRuleForm from '@/components/other/rrule/rrule-form.vue';
import ModalInnerItemMixin from '@/mixins/modal/modal-inner-item';
import { MODALS } from '@/constants';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');
const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');

export default {
  name: MODALS.createPbehavior,

  $_veeValidate: {
    validator: 'new',
  },
  components: { DateTimePicker, RRuleForm },
  mixins: [ModalInnerItemMixin],
  data() {
    const start = new Date();
    const stop = new Date(start.getTime());
    const reasons = ['Problème Habilitation', 'Problème Robot', 'Problème Scénario', 'Autre'];
    const types = ['Pause', 'Maintenance', 'Hors plage horaire de surveillance'];

    return {
      rRuleObject: null,
      form: {
        name: '',
        tstart: start,
        tstop: stop,
        type_: types[0],
        reason: reasons[0],
      },
      selectItems: {
        reasons,
        types,
      },
      serverError: null,
    };
  },
  methods: {
    ...alarmMapActions({ fetchAlarmListWithPreviousParams: 'fetchListWithPreviousParams' }),
    ...pbehaviorMapActions({ createPbehavior: 'create' }),

    changeRRule(value) {
      this.rRuleObject = value;
    },
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.serverError = null;

        const data = {
          ...this.form,

          author: 'Username of current user', // TODO: add this field after login task finish
          filter: {
            _id: { $in: [this.item.d] },
          },
          tstart: this.form.tstart.getTime(),
          tstop: this.form.tstop.getTime(),
        };

        if (this.rRuleObject) {
          data.rrule = this.rRuleObject.toString();
        }

        try {
          await this.createPbehavior(data);
          await this.fetchAlarmListWithPreviousParams();

          this.hideModal();
        } catch (err) {
          if (err.description) {
            this.serverError = err.description;
          }
        }
      }
    },
  },
};
</script>
