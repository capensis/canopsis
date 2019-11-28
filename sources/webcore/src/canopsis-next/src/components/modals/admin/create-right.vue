<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createRight.title') }}
      template(slot="text")
        right-form(v-model="form")
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(:disabled="errors.any()", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { generateRight } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRightMixin from '@/mixins/entities/right';

import RightForm from '@/components/other/right/right-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRight,
  $_veeValidate: {
    validator: 'new',
  },
  components: { RightForm, ModalWrapper },
  mixins: [modalInnerMixin, entitiesRightMixin],
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

          this.$popups.success({ text: this.$t('success.default') });
          this.$modals.hide();
        }
        if (this.config.action) {
          await this.config.action();
        }
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
