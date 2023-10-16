<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <geomap-map-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { mapToForm, formToMap } from '@/helpers/entities/map/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import GeomapMapForm from '@/components/other/map/form/geomap-map-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createGeoMap,
  $_veeValidate: {
    validator: 'new',
  },
  components: { GeomapMapForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: mapToForm(this.modal.config.map),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createGeoMap.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToMap(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
