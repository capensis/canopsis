<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createRight.title') }}
    v-card-text
      right-form(v-model="form")
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { generateRight } from '@/helpers/entities';

import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRightMixin from '@/mixins/entities/right';

import RightForm from '@/components/other/right/right-form.vue';

export default {
  name: MODALS.createRight,
  $_veeValidate: {
    validator: 'new',
  },
  components: { RightForm },
  mixins: [popupMixin, modalInnerMixin, entitiesRightMixin],
  data() {
    return {
      form: {
        _id: '',
        desc: '',
        type: '',
      },
    };
  },
  methods: {
    async submit() {
      try {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          const data = { ...generateRight(), ...this.form };

          await this.createRight({ data });

          this.addSuccessPopup({ text: this.$t('success.default') });
          this.hideModal();
        }
        if (this.config.action) {
          await this.config.action();
        }
      } catch (err) {
        this.addErrorPopup({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
