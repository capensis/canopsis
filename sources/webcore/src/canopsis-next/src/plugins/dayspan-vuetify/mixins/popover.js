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
  methods: {
    openPopover() {
      if (this.isStart) {
        this.menu = true;
      }
    },

    closePopover() {
      this.menu = false;
    },

    triggerClearPlaceholder(menu) {
      if (!menu) {
        setTimeout(() => this.$emit('clear-placeholder'), VUETIFY_ANIMATION_DELAY / 2);
      }
    },
  },
};
