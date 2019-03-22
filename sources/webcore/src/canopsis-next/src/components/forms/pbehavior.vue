<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.createPbehavior.title') }}
      v-card-text
        v-layout(row)
          v-text-field(
          v-model="form.name",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.fields.name')",
          :error-messages="errors.collect('name')",
          name="name"
          )
        v-layout(row)
          date-time-picker-field(
          v-model="form.tstart",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.fields.start')",
          name="tstart"
          )
        v-layout(row)
          date-time-picker-field(
          v-model="form.tstop",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.fields.stop')",
          name="tstop",
          )
        v-layout(v-if="!filter", row)
          v-btn.primary(type="button", @click="showCreateFilterModal") Filter
        r-rule-form(v-model="form.rrule")
        v-layout(row)
          v-combobox(
          v-model="form.reason",
          v-validate="'required'"
          :label="$t('modals.createPbehavior.fields.reason')",
          :items="reasons",
          :error-messages="errors.collect('reason')",
          name="reason",
          )
        v-layout(row)
          v-select(
          v-model="form.type_",
          v-validate="'required'"
          :label="$t('modals.createPbehavior.fields.type')",
          :items="types",
          :error-messages="errors.collect('type')",
          name="type",
          )
        v-layout(row)
          v-textarea(
          :label="$t('modals.createPbehavior.fields.comment')",
          v-model="commentMessage",
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
import { cloneDeep } from 'lodash';

import { MODALS, PAUSE_REASONS, PBEHAVIOR_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import RRuleForm from '@/components/forms/rrule.vue';

/**
 * Modal to create a pbehavior
 */
export default {
  inject: ['$validator'],
  components: { DateTimePickerField, RRuleForm },
  filters: {
    pbehaviorToForm(pbehavior = {}) {
      return {
        name: pbehavior.name || '',
        tstart: pbehavior.tstart ? new Date(pbehavior.tstart * 1000) : new Date(),
        tstop: pbehavior.tstop ? new Date(pbehavior.tstop * 1000) : new Date(),
        filter: cloneDeep(this.filter || {}),
        type_: pbehavior.type_ || '',
        reason: pbehavior.reason || '',
        rrule: pbehavior.rrule || null,
      };
    },
    formToPbehavior() {

    },
  },
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
    value: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      commentMessage: '',
      form: this.$options.filters.pbehaviorToForm(this.value),
    };
  },
  computed: {
    reasons() {
      return Object.values(PAUSE_REASONS);
    },
    types() {
      return Object.values(PBEHAVIOR_TYPES);
    },
  },
  methods: {
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

        if (this.commentMessage !== '') {
          data.comments = [{
            author: this.currentUser.crecord_name,
            message: this.commentMessage,
          }];
        }

        this.$emit('submit', data);
      }
    },
  },
};
</script>
