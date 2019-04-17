<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
        v-btn(icon, dark, @click="hideModal")
          v-icon close
    v-card-text
      pbehavior-comments-form(v-model="config.comments")
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(type="submit", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import PbehaviorCommentsForm from '@/components/other/pbehavior/form/pbehavior-comments-form.vue';

export default {
  name: MODALS.editPbehaviorComment,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorCommentsForm,
  },
  mixins: [modalInnerMixin],
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.config.comments);
        }

        this.hideModal();
      }
    },
  },
};
</script>
