<template>
  <div>
    <groups-side-bar v-if="isShownGroupsSideBar" v-model="isOpen" />
    <top-bar :side-bar-open="isOpen" @toggle:side-bar="toggleSideBar" />
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue';
import { useDisplay } from 'vuetify';

import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import { useAuth } from '@/hooks/auth';

import GroupsSideBar from './partials/groups-side-bar/groups-side-bar.vue';
import TopBar from './partials/top-bar.vue';

const IGNORE_IS_RESIZING_DURATION = 300;
const VuetifyLayoutKey = Symbol.for('vuetify:layout');

defineOptions({
  name: 'TheNavigation',
});

const { currentUser, isShownGroupsTopBar } = useAuth();
const { width } = useDisplay();
const layout = inject(VuetifyLayoutKey, null);

const ignoreLayoutIsResizing = layout?.ignoreIsResizing ?? (() => {});
const isOpen = ref(false);

const isShownGroupsSideBar = computed(() => {
  const groupsNavigationType = currentUser.value?.ui_groups_navigation_type;
  const isSelectedSideBar = groupsNavigationType === GROUPS_NAVIGATION_TYPES.sideBar;
  const isMobileOrTablet = width.value < MEDIA_QUERIES_BREAKPOINTS.t;

  return isSelectedSideBar || isMobileOrTablet || !isShownGroupsTopBar.value;
});

const toggleSideBar = () => {
  ignoreLayoutIsResizing(IGNORE_IS_RESIZING_DURATION);
  isOpen.value = !isOpen.value;
};
</script>
