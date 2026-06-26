<script setup>
import { computed, inject, provide } from 'vue';

import ClickOutside from '@/services/click-outside';

import ModalWrapper from './modal-wrapper.vue';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const modalComponents = inject('$modalComponents', {});

const modalComponent = computed(() => modalComponents[props.modal.name]);

provide('$clickOutside', new ClickOutside());
provide('$modal', props.modal);
</script>

<template>
  <modal-wrapper :modal="modal">
    <component
      :is="modalComponent"
      :modal="modal"
    />
  </modal-wrapper>
</template>
