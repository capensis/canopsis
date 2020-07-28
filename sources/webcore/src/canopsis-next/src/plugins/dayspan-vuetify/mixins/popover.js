import { VUETIFY_ANIMATION_DELAY } from '@/config';

export default {
  data() {
    return {
      isShownPopover: false,
      timer: null,
    };
  },
  watch: {
    menu(menu) {
      if (this.timer) {
        clearTimeout(this.timer);

        this.timer = null;
      }

      if (menu) {
        this.isShownPopover = true;
      } else {
        this.timer = setTimeout(() => this.isShownPopover = false, VUETIFY_ANIMATION_DELAY);
      }
    },
  },
};
