<template>
  <c-modal-wrapper close>
    <template #title>
      {{ t('stateSetting.seeFilterPattern') }}
    </template>
    <template #text>
      <state-setting-inherited-entity-pattern-step disabled />
    </template>
    <template #actions>
      <v-btn :aria-label="t('common.close')" variant="text" @click="closeModal">
        {{ t('common.close') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed } from 'vue';

import { MODALS, PATTERNS_FIELDS } from '@/constants';

import { stateSettingPatternToForm } from '@/helpers/entities/state-setting/form';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import StateSettingInheritedEntityPatternStep from '@/components/other/state-setting/form/steps/state-setting-inherited-entity-pattern-step.vue';

export default {
  name: MODALS.stateSettingInheritedEntityPattern,

  components: {
    StateSettingInheritedEntityPatternStep,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});

    useForm({
      initialValues: {
        inherited_entity_pattern: stateSettingPatternToForm({
          [PATTERNS_FIELDS.entity]: config.value.pattern,
        }),
      },
    });

    const closeModal = () => modals.hide({ id: props.modal.id });

    return {
      t,
      closeModal,
    };
  },
};
</script>
