<script setup>
import { storeToRefs } from 'pinia';

import { usePopupsStore } from '../stores';

import PopupItem from './popup-item.vue';

const popupsStore = usePopupsStore();

const { popups } = storeToRefs(popupsStore);
</script>

<template>
  <teleport to="body">
    <div v-if="popups.length" class="popups">
      <popup-item v-for="popup in popups" v-bind="popup" :key="popup.id" />
    </div>
  </teleport>
</template>

<style lang="scss">
.popups {
  position: fixed;
  z-index: 2147483647;
  right: 2rem;
  top: 75px;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}
</style>
