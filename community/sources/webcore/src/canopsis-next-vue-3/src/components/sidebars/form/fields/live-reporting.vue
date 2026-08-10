<template>
  <settings-button-field
    :title="t('settings.liveReporting.field')"
    :is-empty="isEmpty"
    addable
    removable
    @create="openModal"
    @edit="openModal"
    @delete="clear"
  />
</template>

<script setup>
import { computed } from 'vue';

import { MODALS } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import SettingsButtonField from './button-field.vue';

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();
const modals = useModals();

const isEmpty = computed(() => !props.modelValue?.tstart && !props.modelValue?.tstop);

const openModal = () => {
  modals.show({
    name: MODALS.editLiveReporting,
    config: {
      ...props.modelValue,
      action: value => emit('update:modelValue', value),
    },
  });
};

const clear = () => emit('update:modelValue', {});
</script>
