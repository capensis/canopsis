<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createPbehavior.title') }}
    v-card-text
      pbehavior-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="errors.any()", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';

import {
  commentsToPbehaviorComments,
  formToPbehavior,
  pbehaviorToForm,
  pbehaviorToComments,
  exdatesToPbehaviorExdates,
  pbehaviorToExdates,
} from '@/helpers/forms/pbehavior';

import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorForm,
  },
  mixins: [authMixin, modalInnerMixin],
  data() {
    const { pbehavior = {} } = this.modal.config;

    return {
      form: {
        general: pbehaviorToForm(pbehavior),
        exdate: pbehaviorToExdates(pbehavior),
        comments: pbehaviorToComments(pbehavior),
      },
    };
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const pbehavior = formToPbehavior(this.form.general);

        pbehavior.comments = commentsToPbehaviorComments(this.form.comments);
        pbehavior.exdate = exdatesToPbehaviorExdates(this.form.exdate);

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
