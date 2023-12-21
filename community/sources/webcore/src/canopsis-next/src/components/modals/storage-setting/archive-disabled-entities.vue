<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('storageSetting.entity.archiveDisabled') }}</span>
    </template>
    <template #text="">
      <v-layout column="column">
        <p class="text-subtitle-1 pre-wrap">
          {{ $t('modals.archiveDisabledEntities.text') }}
        </p>
        <v-checkbox
          v-field="form.with_dependencies"
          :label="$t('storageSetting.entity.archiveDependencies')"
          color="primary"
        >
          <template #append="">
            <c-help-icon
              :text="$t('storageSetting.entity.archiveDependenciesHelp')"
              color="info"
              max-width="300"
              top="top"
            />
          </template>
        </v-checkbox>
      </v-layout>
    </template>
    <template #actions="">
      <v-layout
        wrap
        justify-center
      >
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          color="primary"
          @click.prevent="submit"
        >
          {{ $t('common.archive') }}
        </v-btn>
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.archiveDisabledEntities,
  inject: ['$system'],
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      form: {
        with_dependencies: false,
      },
    };
  },
  methods: {
    async submit() {
      await this.config?.action(this.form);

      this.$modals.hide();
    },
  },
};
</script>
