import { writeTextToClipboard } from '@/helpers/clipboard';

export const clipboardMixin = {
  methods: {
    async writeTextToClipboard(text) {
      try {
        await writeTextToClipboard(text);

        this.$popups.success({ text: this.$t('popups.copySuccess') });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('popups.copyError') });
      }
    },
  },
};
