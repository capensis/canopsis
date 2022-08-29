import Observer from '@/services/observer';

export const contextmenuMixin = {
  provide() {
    return {
      $contextmenu: this.$contextmenu,
    };
  },
  props: {
    points: {
      type: Array,
      default: () => [],
    },
  },
  beforeCreate() {
    this.$contextmenu = new Observer();
  },
  methods: {
    handleContextmenu(event) {
      const { x, y } = this.normalizeCursor({ x: event.clientX, y: event.clientY });

      this.$contextmenu.notify({ event, x, y });
    },

    handleShapeContextmenu(shape, event) {
      this.$contextmenu.notify({ event, shape });
    },
  },
};
