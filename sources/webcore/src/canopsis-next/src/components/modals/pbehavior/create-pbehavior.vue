<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createPbehavior.title') }}
    v-card-text
      v-stepper(v-model="stepper", non-linear)
        v-stepper-header
          v-stepper-step(
            :complete="stepper > 1",
            step="1", editable,
            :rules="[() => !hasGeneralFormAnyError]"
          ) {{ $t('modals.createPbehavior.steps.general.title') }}
            small(v-if="hasGeneralFormAnyError") {{ $t('modals.createPbehavior.errors.invalid') }}
          v-divider
          v-stepper-step(
            :complete="stepper > 2",
            step="2",
            editable,
            :rules="[() => !hasFilterEditorAnyError]"
          ) {{ $t('modals.createPbehavior.steps.filter.title') }}
            small(v-if="hasFilterEditorAnyError") {{ $t('modals.createPbehavior.errors.invalid') }}
            small.font-italic.font-weight-light(v-else) {{ $t('common.optional') }}
          v-divider
          v-stepper-step(
            :complete="stepper > 3",
            step="3",
            editable
          ) {{ $t('modals.createPbehavior.steps.rrule.title') }}
            small.font-italic.font-weight-light {{ $t('common.optional') }}
          v-divider
          v-stepper-step(
            :complete="stepper > 4",
            step="4",
            editable
          ) {{ $t('modals.createPbehavior.steps.comments.title') }}
            small.font-italic.font-weight-light {{ $t('common.optional') }}
        v-stepper-items
          v-stepper-content(step="1")
            v-card(flat)
              v-card-text
                pbehavior-general-form(v-model="form", ref="pbehaviorGeneralForm")
          v-stepper-content.pa-0(step="2")
            v-card
              v-card-text
                filter-editor(v-model="form.filter", ref="filterEditor")
          v-stepper-content(step="3")
            v-card
              v-card-text
                r-rule-form(v-model="form.rrule")
                pbehavior-exdates-form(v-if="form.rrule", v-model="exdate")
          v-stepper-content(step="4")
            v-card
              v-card-text
                pbehavior-comments-form(v-model="comments")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="errors.any()", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';
import { cloneDeep, omit, isObject } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';
import convertTimestampToMoment from '@/helpers/date';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';

import PbehaviorGeneralForm from '@/components/other/pbehavior/form/pbehavior-general-form.vue';
import PbehaviorCommentsForm from '@/components/other/pbehavior/form/pbehavior-comments-form.vue';
import PbehaviorExdatesForm from '@/components/other/pbehavior/form/pbehavior-exdates-form.vue';
import RRuleForm from '@/components/forms/rrule.vue';
import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  filters: {
    pbehaviorToForm(pbehavior = {}) {
      let rrule = pbehavior.rrule || null;

      if (pbehavior.rrule && isObject(pbehavior.rrule)) {
        ({ rrule } = pbehavior.rrule);
      }

      return {
        author: pbehavior.author || '',
        name: pbehavior.name || '',
        tstart: pbehavior.tstart ? convertTimestampToMoment(pbehavior.tstart).toDate() : new Date(),
        tstop: pbehavior.tstop ? convertTimestampToMoment(pbehavior.tstop).toDate() : moment().add(1, 'd').toDate(),
        timezone: 'Europe/Paris',
        enabled: pbehavior.enabled || true,
        filter: cloneDeep(pbehavior.filter || {}),
        type_: pbehavior.type_ || '',
        reason: pbehavior.reason || '',
        rrule,
        connector: pbehavior.connector || '',
        connector_name: pbehavior.connector_name || '',
      };
    },

    pbehaviorToComments(pbehavior = {}) {
      const comments = pbehavior.comments || [];

      return comments.map(comment => ({
        ...comment,

        key: uid(),
      }));
    },

    pbehaviorToExdate(pbehavior = {}) {
      const exdate = pbehavior.exdate || [];

      return exdate.map(unix => ({
        value: new Date(unix * 1000),
        key: uid(),
      }));
    },

    formToPbehavior(form) {
      return {
        ...form,
        comments: [],
        tstart: moment(form.tstart).unix(),
        tstop: moment(form.tstop).unix(),
      };
    },

    commentsToPbehaviorComments(comments) {
      return comments.map(comment => omit(comment, ['key', 'ts']));
    },

    exdateToPbehaviorExdate(exdate) {
      return exdate.filter(({ value }) => value).map(({ value }) => moment(value).unix());
    },
  },
  components: {
    PbehaviorGeneralForm,
    PbehaviorCommentsForm,
    PbehaviorExdatesForm,
    RRuleForm,
    FilterEditor,
  },
  mixins: [authMixin, modalInnerMixin],
  data() {
    const { pbehavior = {} } = this.modal.config;

    return {
      stepper: 1,
      hasGeneralFormAnyError: false,
      hasFilterEditorAnyError: false,
      form: this.$options.filters.pbehaviorToForm(pbehavior),
      exdate: this.$options.filters.pbehaviorToExdate(pbehavior),
      comments: this.$options.filters.pbehaviorToComments(pbehavior),
    };
  },
  mounted() {
    this.$watch(() => this.$refs.pbehaviorGeneralForm.hasAnyError, (value) => {
      this.hasGeneralFormAnyError = value;
    });

    this.$watch(() => this.$refs.filterEditor.hasAnyError, (value) => {
      this.hasFilterEditorAnyError = value;
    });
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const pbehavior = this.$options.filters.formToPbehavior(this.form);

        pbehavior.comments = this.$options.filters.commentsToPbehaviorComments(this.comments);
        if (pbehavior.rrule) {
          pbehavior.exdate = this.$options.filters.exdateToPbehaviorExdate(this.exdate);
        }

        if (!pbehavior.author) {
          pbehavior.author = this.currentUser._id;
        }

        if (this.config.action) {
          await this.config.action(pbehavior);
        }

        this.hideModal();
      }
    },
  },
};
</script>
