import { TOP_BAR_HEIGHT } from '@/config';
import { ALARMS_LIST_HEADER_OPACITY_DELAY } from '@/constants';

export const widgetStickyAlarmMixin = {
  props: {
    stickyHeader: {
      type: Boolean,
      default: false,
    },
    stickyHorizontalScroll: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    topBarHeight() {
      return this.shownHeader ? TOP_BAR_HEIGHT : 0;
    },

    tableWrapper() {
      return this.$el.querySelector('.alarms-list-table .v-data-table__wrapper');
    },

    tableHeader() {
      return this.$el.querySelector('.alarms-list-table .v-data-table-header');
    },

    tableBody() {
      return this.$el.querySelector('.alarms-list-table table  tbody');
    },

    someSticky() {
      return this.stickyHeader || this.stickyHorizontalScroll;
    },
  },

  watch: {
    someSticky(someSticky) {
      if (someSticky) {
        this.calculateStickyPositions();
        this.setHeaderPosition();
        this.setHorizontalScrollPosition();

        if (this.stickyHorizontalScroll) {
          this.connectHorizontalScroll();
        }

        window.addEventListener('scroll', this.stickyScrollHandler);
      } else {
        window.removeEventListener('scroll', this.stickyScrollHandler);

        this.resetHeaderPosition();
      }
    },
  },

  created() {
    this.actionsTranslateY = 0;
    this.translateY = 0;
    this.previousTranslateY = 0;
    this.previousScrollTranslateY = 0;
  },

  async mounted() {
    if (this.someSticky) {
      window.addEventListener('scroll', this.stickyScrollHandler);

      if (this.stickyHorizontalScroll) {
        this.connectHorizontalScroll();
      }
    }
  },

  beforeDestroy() {
    window.removeEventListener('scroll', this.stickyScrollHandler);
  },

  methods: {
    connectHorizontalScroll() {
      this.$refs.horizontalScrollbar?.connect(this.tableWrapper);
    },

    startScrollingForHeader() {
      if (this.translateY !== this.previousTranslateY) {
        this.tableHeader.style.opacity = '0';

        if (this.$refs.actions) {
          this.$refs.actions.style.opacity = '0';
        }
      }
    },

    startScrollingForScrollbar() {
      if (this.$refs.horizontalScrollbar && this.previousScrollTranslateY !== this.scrollBarTranslateY) {
        this.$refs.horizontalScrollbar.$el.style.opacity = '0';
      }
    },

    startScrolling() {
      this.scrolling = true;
    },

    finishScrollingForHeader() {
      if (!Number(this.tableHeader.style.opacity)) {
        this.tableHeader.style.opacity = '1.0';

        if (this.$refs.actions) {
          this.$refs.actions.style.opacity = '1.0';
        }
      }
    },

    finishScrollingForScrollbar() {
      if (this.$refs.horizontalScrollbar?.$el?.style?.opacity === '0') {
        this.$refs.horizontalScrollbar.$el.style.opacity = '1';
      }
    },

    finishScrolling() {
      this.scrolling = false;
    },

    fullFinishScrolling() {
      this.finishScrolling();
      this.finishScrollingForHeader();
      this.finishScrollingForScrollbar();
    },

    clearFinishTimer() {
      if (this.finishTimer) {
        clearTimeout(this.finishTimer);
      }
    },

    setHeaderPosition() {
      if (!this.translateY || !this.actionsTranslateY) {
        this.removeShadowFromHeader();
      } else {
        this.addShadowToHeader();
      }

      this.tableHeader.style.transform = `translateY(${this.translateY}px)`;

      if (this.$refs.actions) {
        this.$refs.actions.style.transform = `translateY(${this.actionsTranslateY}px)`;
      }

      if (this.translateY === 0) {
        this.finishScrollingForHeader();
      }
    },

    setHorizontalScrollPosition() {
      this.$refs.horizontalScrollbar?.setTranslateY?.(this.scrollBarTranslateY);

      if (this.scrollBarTranslateY === 0) {
        this.finishScrollingForScrollbar();
      }
    },

    calculateHeaderPosition() {
      const { height: bodyHeight } = this.tableBody.getBoundingClientRect();
      const { top: headerTop } = this.tableHeader.getBoundingClientRect();
      const { top: actionsTop = 0, height: actionsHeight = 0 } = this.$refs.actions?.getBoundingClientRect() ?? {};

      const offset = headerTop - this.translateY - this.topBarHeight - actionsHeight;
      const actionsOffset = actionsTop - this.actionsTranslateY - this.topBarHeight;

      this.previousTranslateY = this.actionsTranslateY;
      this.translateY = Math.min(bodyHeight, Math.max(0, -offset));
      this.actionsTranslateY = Math.min(bodyHeight, Math.max(0, -actionsOffset));
    },

    calculateScrollPosition() {
      const { bottom: bodyBottom } = this.tableBody.getBoundingClientRect();
      const {
        height: scrollBarHeight = 0,
      } = this.$refs.horizontalScrollbar?.$el.getBoundingClientRect() ?? {};

      this.previousScrollTranslateY = this.scrollBarTranslateY;
      this.scrollBarTranslateY = Math.min(window.innerHeight - bodyBottom - scrollBarHeight, 0);
    },

    setHorizontalScrollbarWidth() {
      this.$refs.horizontalScrollbar.setContentWidth(this.tableWrapper?.scrollWidth ?? 0);
    },

    calculateStickyPositions() {
      if (this.stickyHeader) {
        this.calculateHeaderPosition();
      }

      if (this.stickyHorizontalScroll) {
        this.calculateScrollPosition();
      }
    },

    addShadowToHeader() {
      this.tableHeader.classList.add('head-shadow');
    },

    removeShadowFromHeader() {
      this.tableHeader.classList.remove('head-shadow');
    },

    stickyScrollHandler() {
      this.clearFinishTimer();

      this.calculateStickyPositions();
      this.setHeaderPosition();

      if (this.stickyHorizontalScroll) {
        this.setHorizontalScrollbarWidth();
        this.setHorizontalScrollPosition();
      }

      if (!this.actionsTranslateY || !this.translateY) {
        this.finishScrollingForHeader();

        if (!this.scrollBarTranslateY) {
          this.finishScrollingForScrollbar();
          this.finishScrolling();
          return;
        }
      }

      if (!this.scrolling) {
        this.startScrolling();
        this.startScrollingForHeader();
        this.startScrollingForScrollbar();
      }

      this.finishTimer = setTimeout(this.fullFinishScrolling, ALARMS_LIST_HEADER_OPACITY_DELAY);
    },

    resetHeaderPosition() {
      this.translateY = 0;
      this.actionsTranslateY = 0;
      this.previousTranslateY = 0;
      this.previousScrollTranslateY = 0;

      this.setHeaderPosition();
      this.clearFinishTimer();
      this.removeShadowFromHeader();
    },
  },
};
