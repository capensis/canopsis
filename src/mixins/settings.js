import Settings from '@/components/other/settings/settings.vue';

export default {
  components: { Settings },
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
