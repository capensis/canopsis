<template lang="pug">
  v-form(@submit.prevent="submit")
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
        v-layout(v-if="!filter", row)
          v-btn.primary(type="button", @click="showCreateFilterModal") Filter
        r-rule-form(@input="changeRRule")
        v-layout(row)
          v-select(
          label="Reason",
          v-model="form.reason",
          :items="selectItems.reasons",
          :error-messages="errors.collect('reason')",
          name="reason",
          v-validate="'required'"
          )
        v-layout(row)
          v-select(
          label="Type",
          v-model="form.type_",
          :items="selectItems.types",
          :error-messages="errors.collect('type')",
          name="type",
          v-validate="'required'"
          )
        v-layout(row)
          v-alert(:value="serverError", type="error")
            span {{ serverError }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="cancel", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';

import { MODALS } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';

import DateTimePicker from '@/components/forms/fields/date-time-picker.vue';
import RRuleForm from '@/components/forms/rrule.vue';

/**
 * Modal to create a pbehavior
 */
export default {
  inject: ['$validator'],
  components: { DateTimePicker, RRuleForm },
  mixins: [authMixin, modalMixin],
  props: {
    serverError: {
      type: String,
      default: null,
    },
    filter: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      rRuleObject: null,
      form: {
        name: '',
        tstart: new Date(),
        tstop: new Date(),
        filter: this.filter,
        type_: '',
        reason: '',
      },
      selectItems: {
        reasons: ['Problème Habilitation', 'Problème Robot', 'Problème Scénario', 'Autre'],
        types: ['Pause', 'Maintenance', 'Hors plage horaire de surveillance'],
      },
    };
  },
  methods: {
    changeRRule(value) {
      this.rRuleObject = value;
    },

    showCreateFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'Pbehavior filter',
          hiddenFields: ['title'],
          filter: {
            filter: this.form.filter || {},
          },
          action: ({ filter }) => this.form.filter = filter,
        },
      });
    },

    cancel() {
      this.$emit('cancel');
    },

    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const data = {
          ...this.form,

          author: this.currentUser.crecord_name,
          tstart: moment(this.form.tstart).unix(),
          tstop: moment(this.form.tstop).unix(),
        };

        if (this.rRuleObject) {
          data.rrule = this.rRuleObject.toString();
        }

        this.$emit('submit', data);
      }
    },
  },
};
</script>
