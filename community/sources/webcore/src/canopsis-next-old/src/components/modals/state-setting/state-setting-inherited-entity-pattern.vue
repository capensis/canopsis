<template>
  <modal-wrapper close>
    <template #title="">
      {{ $t('stateSetting.seeFilterPattern') }}
    </template>
    <template #text="">
      <state-setting-inherited-entity-pattern-step :patterns="pattern" disabled />
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { stateSettingPatternToForm } from '@/helpers/entities/state-setting/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { authMixin } from '@/mixins/auth';

import StateSettingInheritedEntityPatternStep from '@/components/other/state-setting/form/steps/state-setting-inherited-entity-pattern-step.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.stateSettingInheritedEntityPattern,
  components: {
    StateSettingInheritedEntityPatternStep,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    authMixin,
  ],
  computed: {
    pattern() {
      return stateSettingPatternToForm(this.config.pattern);
    },
  },
};
</script>
