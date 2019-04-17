<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createPbehavior.title') }}
    v-card-text
      pbehavior-form(v-model="form")
      pbehavior-exdates-form.mt-2(v-show="form.rrule", v-model="exdate")
      pbehavior-comments-form.mt-2(v-model="comments")
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

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';

import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';
import PbehaviorExdatesForm from '@/components/other/pbehavior/form/pbehavior-exdates-form.vue';
import PbehaviorCommentsForm from '@/components/other/pbehavior/form/pbehavior-comments-form.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  filters: {
    pbehaviorToForm(pbehavior = {}) {
      let rrule = pbehavior.rrule || '';

      if (pbehavior.rrule && isObject(pbehavior.rrule)) {
        ({ rrule } = pbehavior.rrule);
      }

      return {
        author: pbehavior.author || '',
        name: pbehavior.name || '',
        tstart: pbehavior.tstart ? new Date(pbehavior.tstart * 1000) : new Date(),
        tstop: pbehavior.tstop ? new Date(pbehavior.tstop * 1000) : new Date(),
        filter: cloneDeep(pbehavior.filter || {}),
        type_: pbehavior.type_ || '',
        reason: pbehavior.reason || '',
        rrule,
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
      return comments.map(comment => omit(comment, ['key']));
    },

    exdateToPbehaviorExdate(exdate) {
      return exdate.filter(({ value }) => value).map(({ value }) => moment(value).unix());
    },
  },
  components: {
    PbehaviorForm,
    PbehaviorExdatesForm,
    PbehaviorCommentsForm,
  },
  mixins: [authMixin, modalInnerMixin],
  data() {
    const { pbehavior = {} } = this.modal.config;

    return {
      form: this.$options.filters.pbehaviorToForm(pbehavior),
      exdate: this.$options.filters.pbehaviorToExdate(pbehavior),
      comments: this.$options.filters.pbehaviorToComments(pbehavior),
    };
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const pbehavior = this.$options.filters.formToPbehavior(this.form);

        pbehavior.comments = this.$options.filters.commentsToPbehaviorComments(this.comments);
        pbehavior.exdate = this.$options.filters.exdateToPbehaviorExdate(this.exdate);

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
