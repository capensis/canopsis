import { VUETIFY_ANIMATION_DELAY } from '@/config';

export default {
  inject: ['$dayspanOptions'],
  props: {
    popoverProps: {
      validate(x) {
        return this.$dsValidate(x, 'popoverProps');
      },
      default() {
        const defaultPopoverOptions = this.$dsDefaults().popoverProps || {};
        const calendarPopoverOptions = this.$dayspanOptions.getOptions(this.$options.name).popoverProps || {};

        return {
          ...defaultPopoverOptions,
          ...calendarPopoverOptions,
        };
      },
    },
  },
  computed: {
    menuDisabled() {
      return !this.hasPopover || this.placeholder.data.resizing || this.placeholder.data.moving;
    },
  },
  data() {
    return {
      menu: false,
      isShownPopover: false,
      openTimer: null,
      closeTimer: null,
    };
  },
  watch: {
    menu(menu) {
      if (this.openTimer) {
        clearTimeout(this.openTimer);

        this.openTimer = null;
      }

      if (menu) {
        this.isShownPopover = true;
      } else {
        this.openTimer = setTimeout(() => this.isShownPopover = false, VUETIFY_ANIMATION_DELAY);
      }
    },
  },
  methods: {
    openPopover(menu = true) {
      if (this.isStart) {
        this.menu = menu;
      }
    },

    closePopover() {
      this.menu = false;
    },

    triggerClearPlaceholder(menu) {
      if (this.closeTimer) {
        clearTimeout(this.closeTimer);

        this.closeTimer = null;
      }

      if (!menu) {
        this.closeTimer = setTimeout(() => this.$emit('clear-placeholder'), VUETIFY_ANIMATION_DELAY / 2);
      }
    },
  },
};
