<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-text-field
          v-model="form.title"
          v-validate="'required'"
          :label="$t('common.title')"
          :error-messages="errors.collect('title')"
          name="title"
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
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
        <v-tooltip
          v-if="config.deletable"
          :disabled="group.deletable"
          top
        >
          <template #activator="{ on }">
            <div
              class="ml-2"
              v-on="on"
            >
              <v-btn
                :disabled="submitting || !group.deletable"
                :outlined="$system.dark"
                class="error"
                color="error"
                @click="remove"
              >
                {{ $t('common.delete') }}
              </v-btn>
            </div>
          </template>
          <span>{{ $t('modals.group.errors.isNotEmpty') }}</span>
        </v-tooltip>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { groupToRequest } from '@/helpers/entities/view/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createGroup,
  $_veeValidate: {
    validator: 'new',
  },
  inject: ['$system'],
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        title: this.modal.config.group.title || '',
      },
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.group.create.title');
    },

    group() {
      return this.config.group;
    },
  },
  methods: {
    remove() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.config.remove?.();

              this.$modals.hide();
            } catch (err) {
              this.$popups.error({ text: this.$t('modals.group.errors.isNotEmpty') });
            }
          },
        },
      });
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.config.action?.(groupToRequest({ ...this.group, ...this.form }));

        this.$modals.hide();
      }
    },
  },
};
</script>
