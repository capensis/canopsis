export const contextmenuMixin = {
  data() {
    return {
      shownMenu: false,
      pageX: 0,
      pageY: 0,
      clientX: 0,
      clientY: 0,
      pointData: undefined,
    };
  },
  computed: {
    contextmenuItems() {
      return [
        {
          text: this.$t('map.addPoint'),
          action: this.addPoint,
        },
      ];
    },
  },
  methods: {
    openContextmenu() {
      this.shownMenu = true;
    },

    closeContextmenu() {
      this.shownMenu = false;
    },

    setPageOffsetByEvent(event) {
      this.pageX = event.pageX;
      this.pageY = event.pageY - window.scrollY;
      this.clientX = event.clientX;
      this.clientY = event.clientY;
    },

    handleContextmenu(event) {
      if (this.shownMenu) {
        return;
      }

      this.setPageOffsetByEvent(event);
      this.openContextmenu();
    },

    addPoint() {
      const { x, y } = this.normalizeCursor({ x: this.clientX, y: this.clientY });

      this.pointData = { x, y };
    },
  },
};
