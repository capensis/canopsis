<template lang="pug">
  div
    side-bar-wrapper(:navigationDrawerProps="navigationDrawerProps[sideBarName] || defaultNavigationDrawerProps")
      component(:is="sideBarName", :config="sideBarConfig")
</template>

<script>
import sideBarInnerMixin from '@/mixins/side-bar/side-bar-inner';

import SideBarWrapper from './side-bar-wrapper.vue';
import AlarmSettingsFields from './settings/widgets/alarm.vue';
import ContextSettingsFields from './settings/widgets/context.vue';

/**
 * Wrapper for all side-bars
 */
export default {
  components: {
    SideBarWrapper,
    AlarmSettingsFields,
    ContextSettingsFields,
  },
  mixins: [sideBarInnerMixin],
  data() {
    return {
      navigationDrawerProps: {},
      defaultNavigationDrawerProps: {
        right: true,
        fixed: true,
        temporary: true,
        width: 400,
      },
    };
  },
};
</script>
