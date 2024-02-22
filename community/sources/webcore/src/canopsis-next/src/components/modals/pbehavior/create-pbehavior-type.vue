<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createPbehaviorType.title') }}</span>
      </template>
      <template #text="">
        <pbehavior-type-form
          v-model="form"
          :only-color="onlyColor"
          :pending-priority="pendingPriority"
        />
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
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorTypeToForm } from '@/helpers/entities/pbehavior/type/form';

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
      pendingPriority: false,
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
      this.pendingPriority = true;

      try {
        const { priority } = await this.fetchNextPbehaviorTypePriority();

        this.form.priority = priority;
      } catch (err) {
        console.error(err);
      } finally {
        this.pendingPriority = false;
      }
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
