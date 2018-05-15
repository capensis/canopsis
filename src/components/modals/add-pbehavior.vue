<template lang="pug">
  v-form(@submit.prevent="submit", slot-scope="slotProps")
    v-card
      v-card-title
        span.headline {{ $t('modals.addPbehavior.title') }}
      v-card-text
        v-layout(row)
          v-text-field(
          :label="$t('modals.addPbehavior.fields.name')",
          :error-messages="errors.collect('name')",
          v-model="form.name",
          v-validate="'required'",
          data-vv-name="name"
          )
        v-layout(row)
          date-time-picker(
          :label="$t('modals.addPbehavior.fields.start')",
          v-model="form.tstart",
          name="tstart",
          rules="required",
          )
        v-layout(row)
          date-time-picker(
          :label="$t('modals.addPbehavior.fields.stop')",
          v-model="form.tstop",
          name="tstop",
          rules="required"
          )
        r-rule-form(:tstart="form.tstart", @input="changeRRule")
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
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import DateTimePicker from '@/components/forms/date-time-picker.vue';
import RRuleForm from '@/components/other/rrule/rrule-form.vue';

const { mapGetters: modalMapGetters, mapActions: modalMapActions } = createNamespacedHelpers('modal');
const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('entities/pbehavior');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { DateTimePicker, RRuleForm },
  data() {
    const now = new Date();
    const reasons = ['Problème Habilitation', 'Problème Robot', 'Problème Scénario', 'Autre'];
    const types = ['Pause', 'Maintenance', 'Hors plage horaire de surveillance'];

    return {
      rRuleObject: null,
      form: {
        name: '',
        tstart: now,
        tstop: now,
        type_: types[0],
        reason: reasons[0],
      },
      selectItems: {
        reasons,
        types,
      },
    };
  },
  computed: {
    ...modalMapGetters(['config']),
  },
  methods: {
    ...modalMapActions({ hideModal: 'hide' }),

    ...pbehaviorMapActions({ createPbehavior: 'create' }),

    changeRRule(value) {
      this.rRuleObject = value;
    },
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const data = {
          ...this.form,

          author: 'Username of current user',
          filter: {
            _id: {
              $in: [this.config.d],
            },
          },
          tstart: this.form.tstart.getTime(),
          tstop: this.form.tstop.getTime(),
        };

        if (this.rRuleObject) {
          data.rrule = this.rRuleObject.toString();
        }

        await this.createPbehavior(data);

        this.hideModal();
      }
    },
  },
};
</script>
