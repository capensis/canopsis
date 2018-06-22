// COMPONENTS
import SettingsWrapper from '@/components/other/settings/settings-wrapper.vue';

export default {
  components: { SettingsWrapper },
  data() {
    return {
      isSettingsOpen: false,
    };
  },
  methods: {
    openSettings() {
      this.isSettingsOpen = true;
    },
    closeSettings() {
      this.isSettingsOpen = false;
    },
  },
};
