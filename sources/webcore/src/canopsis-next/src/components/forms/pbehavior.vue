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
          v-validate="tstopRules",
          :label="$t('modals.createPbehavior.fields.stop')",
          name="tstop"
          )
        v-layout(v-if="!filter", row)
          v-btn.primary(type="button", @click="showCreateFilterModal") Filter
        r-rule-form(v-model="form.rrule")
        v-layout(row)
          v-combobox(
          v-model="form.reason",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.fields.reason')",
          :items="reasons",
          :error-messages="errors.collect('reason')",
          name="reason",
          )
        v-layout(row)
          v-select(
          v-model="form.type_",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.fields.type')",
          :items="types",
          :error-messages="errors.collect('type')",
          name="type",
          )
        v-layout(row)
          strong Comments
        v-layout(v-for="(comment, key) in comments", :key="key", row, wrap, allign-center)
          v-flex(xs11)
            v-textarea(
            :disabled="!!comment._id",
            :label="$t('modals.createPbehavior.fields.message')",
            :value="comment.message",
            @input="updateCommentMessage(key, $event)"
            )
          v-flex(xs1)
            v-btn(color="error", icon, @click="removeComment(key)")
              v-icon delete
        v-layout(row)
          v-btn.primary(type="button", @click="addComment") Add comment
        v-layout(row)
          v-alert(:value="errors.has('server')", type="error")
            span {{ errors.first('server') }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="cancel", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';
import { cloneDeep } from 'lodash';

import { MODALS, PAUSE_REASONS, PBEHAVIOR_TYPES } from '@/constants';

import uid from '@/helpers/uid';

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
        filter: cloneDeep(pbehavior.filter || {}),
        type_: pbehavior.type_ || '',
        reason: pbehavior.reason || '',
        rrule: pbehavior.rrule || null,
      };
    },

    pbehaviorToComments(pbehavior = {}) {
      const comments = pbehavior.comments || [];

      return comments.reduce((acc, comment) => {
        acc[uid()] = comment;

        return acc;
      }, {});
    },

    formToPbehavior(form) {
      return {
        ...form,

        comments: [],
        tstart: moment(form.tstart).unix(),
        tstop: moment(form.tstop).unix(),
      };
    },
  },
  mixins: [authMixin, modalMixin],
  model: {
    prop: 'pbehavior',
    event: 'input',
  },
  props: {
    serverError: {
      type: String,
      default: null,
    },
    filter: {
      type: Object,
      default: null,
    },
    pbehavior: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      form: this.$options.filters.pbehaviorToForm(this.pbehavior),
      comments: this.$options.filters.pbehaviorToComments(this.pbehavior),
    };
  },
  computed: {
    reasons() {
      return Object.values(PAUSE_REASONS);
    },

    types() {
      return Object.values(PBEHAVIOR_TYPES);
    },

    tstopRules() {
      const rules = { required: true };

      if (this.form.tstart) {
        rules.after = [this.form.tstart];
      }

      return rules;
    },
  },
  methods: {
    removeComment(key) {
      this.$delete(this.comments, key);
    },

    updateCommentMessage(key, value) {
      this.$set(this.comments[key], 'message', value);
    },

    addComment() {
      this.$set(this.comments, uid(), {
        message: '',
        author: this.currentUser.crecord_name,
      });
    },

    showCreateFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
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
        const pbehavior = this.$options.filters.formToPbehavior(this.form);

        pbehavior.author = this.currentUser.crecord_name;
        pbehavior.comments = Object.values(this.comments);

        this.$emit('submit', pbehavior);
      }
    },
  },
};
</script>
