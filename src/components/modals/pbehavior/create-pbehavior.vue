<template lang="pug">
  v-form(@submit.prevent="submit", slot-scope="slotProps")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
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
          v-alert(:value="serverError", type="error")
            span {{ serverError }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import DateTimePicker from '@/components/forms/date-time-picker.vue';
import RRuleForm from '@/components/forms/rrule.vue';
import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import authMixin from '@/mixins/auth';
import { MODALS } from '@/constants';

const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');

/**
 * Modal to create a pbehavior
 */
export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DateTimePicker, RRuleForm },
  mixins: [modalInnerItemsMixin, authMixin],
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

          author: this.currentUser.crecord_name,
          filter: {
            _id: { $in: this.items.map(v => v._id) },
          },
          tstart: this.form.tstart.getTime(),
          tstop: this.form.tstop.getTime(),
        };

        if (this.rRuleObject) {
          data.rrule = this.rRuleObject.toString();
        }

        try {
          await this.createPbehavior({ data, parents: this.items, parentsType: this.config.itemsType });

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
