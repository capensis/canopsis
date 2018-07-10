import SettingsWrapper from '@/components/other/settings/settings-wrapper.vue';

/**
 * @mixin For the setting bar's displaying
 */
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
