<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createPbehaviorType.title') }}
      template(#text="")
        pbehavior-type-form(v-model="form", :only-color="onlyColor")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorTypeToForm } from '@/helpers/forms/type-pbehavior';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { entitiesPbehaviorTypeMixin } from '@/mixins/entities/pbehavior/types';

import PbehaviorTypeForm from '@/components/other/pbehavior/types/form/pbehavior-type-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorType,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
    PbehaviorTypeForm,
  },
  mixins: [
    modalInnerMixin,
    entitiesPbehaviorTypeMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      form: pbehaviorTypeToForm(this.modal.config.pbehaviorType),
    };
  },
  computed: {
    pbehaviorType() {
      return this.modal.config.pbehaviorType;
    },

    onlyColor() {
      return this.pbehaviorType?.default;
    },

    isNew() {
      return !this.pbehaviorType?._id;
    },
  },
  mounted() {
    if (this.isNew) {
      this.setMinimalPriority();
    }
  },
  methods: {
    async setMinimalPriority() {
      const { priority } = await this.fetchNextPbehaviorTypePriority();

      this.form.priority = priority;
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
