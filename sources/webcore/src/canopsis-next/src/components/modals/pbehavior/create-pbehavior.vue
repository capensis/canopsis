<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createPbehavior.title') }}
      template(slot="text")
        pbehavior-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import {
  commentsToPbehaviorComments,
  formToPbehavior,
  pbehaviorToForm,
  pbehaviorToComments,
  exdatesToPbehaviorExdates,
  pbehaviorToExdates,
} from '@/helpers/forms/pbehavior';

import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorForm, ModalWrapper },
  mixins: [authMixin, modalInnerMixin, submittableMixin()],
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

        this.$modals.hide();
      }
    },
  },
};
</script>
