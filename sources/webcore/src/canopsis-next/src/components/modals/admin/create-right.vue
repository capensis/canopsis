<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createRight.title') }}
      template(slot="text")
        right-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { generateRight } from '@/helpers/entities';

import entitiesRightMixin from '@/mixins/entities/right';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import RightForm from '@/components/other/right/form/right-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRight,
  $_veeValidate: {
    validator: 'new',
  },
  components: { RightForm, ModalWrapper },
  mixins: [
    entitiesRightMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
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
    },
  },
};
</script>
