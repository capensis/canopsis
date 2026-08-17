import { TOP_BAR_HEIGHT } from '@/config';
import { ALARMS_LIST_HEADER_OPACITY_DELAY } from '@/constants';

import { shownHeaderMixin } from '@/mixins/shown-header';

export const widgetStickyAlarmMixin = {
  mixins: [shownHeaderMixin],

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

  data() {
    return {
      actionsTranslateY: 0,
      previousScrollTranslateY: 0,
      previousTranslateY: 0,
      scrollBarTranslateY: 0,
      scrolling: false,
      translateY: 0,
    };
  },

  computed: {
    topBarHeight() {
      return this.shownHeader ? TOP_BAR_HEIGHT : 0;
    },

    stickyTableWrapper() {
      return this.$refs.dataTable?.$el?.querySelector('.v-table__wrapper');
    },

    stickyTable() {
      return this.$refs.dataTable?.$el?.querySelector('table');
    },

    stickyTableHeader() {
      return this.$refs.dataTable?.$el?.querySelector('thead');
    },

    stickyTableBody() {
      return this.$refs.dataTable?.$el?.querySelector('tbody');
    },

    stickyFlagsArray() {
      return [this.stickyHeader, this.stickyHorizontalScroll];
    },

    someSticky() {
      return this.stickyFlagsArray.some(Boolean);
    },
  },

  watch: {
    stickyFlagsArray() {
      if (this.someSticky) {
        this.scheduleStickyScrollHandler();
        this.registerScrollHandler();
      } else {
        this.unregisterScrollHandler();
        this.resetHeaderPosition();
        this.$refs.horizontalScrollbar?.disconnect?.();
      }
    },
  },

  mounted() {
    if (this.someSticky) {
      this.registerScrollHandler();
      this.scheduleStickyScrollHandler();
    }
  },

  beforeUnmount() {
    this.unregisterScrollHandler();
    this.clearFinishTimer();
  },

  methods: {
    registerScrollHandler() {
      if (this.stickyHorizontalScroll) {
        this.connectHorizontalScroll();
      }

      if (this.registeredScrollHandler) {
        return;
      }

      window.addEventListener('scroll', this.stickyScrollHandler);
      window.addEventListener('resize', this.stickyScrollHandler);
      this.registeredScrollHandler = true;
    },

    unregisterScrollHandler() {
      window.removeEventListener('scroll', this.stickyScrollHandler);
      window.removeEventListener('resize', this.stickyScrollHandler);
      this.registeredScrollHandler = false;
    },

    connectHorizontalScroll() {
      this.$nextTick(() => this.$refs.horizontalScrollbar?.connect(this.stickyTableWrapper));
    },

    getStickyActionsEl() {
      const actions = this.$refs.actions;

      if (actions instanceof HTMLElement) {
        return actions;
      }

      if (actions?.$el instanceof HTMLElement) {
        return actions.$el;
      }

      return this.$el?.querySelector('.alarms-list-table__top-pagination') ?? null;
    },

    startScrollingForHeader() {
      if (this.stickyTableHeader && this.translateY !== this.previousTranslateY) {
        this.stickyTableHeader.style.opacity = '0';

        const stickyActions = this.getStickyActionsEl();

        if (stickyActions) {
          stickyActions.style.opacity = '0';
        }
      }
    },

    startScrollingForScrollbar() {
      if (this.$refs.horizontalScrollbar && this.previousScrollTranslateY !== this.scrollBarTranslateY) {
        this.$refs.horizontalScrollbar.setOpacity('0');
      }
    },

    startScrolling() {
      this.scrolling = true;
    },

    finishScrollingForHeader() {
      if (this.stickyTableHeader && !Number(this.stickyTableHeader.style.opacity)) {
        this.stickyTableHeader.style.opacity = '1.0';

        const stickyActions = this.getStickyActionsEl();

        if (stickyActions) {
          stickyActions.style.opacity = '1.0';
        }
      }
    },

    finishScrollingForScrollbar() {
      if (this.$refs.horizontalScrollbar?.$el?.style?.opacity === '0') {
        this.$refs.horizontalScrollbar.setOpacity('1');
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
      if (!this.stickyHeader || !this.stickyTableHeader) {
        return;
      }

      if (!this.translateY) {
        this.removeShadowFromHeader();
      } else {
        this.addShadowToHeader();
      }

      this.stickyTableHeader.style.transform = `translateY(${this.translateY}px)`;

      const stickyActions = this.getStickyActionsEl();

      if (stickyActions) {
        stickyActions.style.transform = `translateY(${this.actionsTranslateY}px)`;
      }

      if (this.translateY === 0) {
        this.finishScrollingForHeader();
      }
    },

    setHeaderWidth() {
      if (!this.stickyTableHeader || !this.stickyTable) {
        return;
      }

      this.stickyTableHeader.style.width = `${this.stickyTable.scrollWidth}px`;
    },

    setHorizontalScrollPosition() {
      this.$refs.horizontalScrollbar?.setTranslateY?.(this.scrollBarTranslateY);

      if (this.scrollBarTranslateY === 0) {
        this.finishScrollingForScrollbar();
      }
    },

    calculateHeaderPosition() {
      if (!this.stickyTableBody || !this.stickyTableHeader) {
        return;
      }

      const { height: bodyHeight } = this.stickyTableBody.getBoundingClientRect();
      const { top: headerTop } = this.stickyTableHeader.getBoundingClientRect();
      const {
        top: actionsTop = 0,
        height: actionsHeight = 0,
      } = this.getStickyActionsEl()?.getBoundingClientRect() ?? {};
      const offset = headerTop - this.translateY - this.topBarHeight - actionsHeight;
      const actionsOffset = actionsTop - this.actionsTranslateY - this.topBarHeight;

      this.previousTranslateY = this.actionsTranslateY;
      this.translateY = Math.min(bodyHeight, Math.max(0, -offset));
      this.actionsTranslateY = Math.min(bodyHeight, Math.max(0, -actionsOffset));
    },

    calculateScrollPosition() {
      if (!this.stickyTableBody) {
        return;
      }

      const { bottom: bodyBottom } = this.stickyTableBody.getBoundingClientRect();
      const {
        height: scrollBarHeight = 0,
      } = this.$refs.horizontalScrollbar?.$el.getBoundingClientRect() ?? {};

      this.previousScrollTranslateY = this.scrollBarTranslateY;
      this.scrollBarTranslateY = Math.min(window.innerHeight - bodyBottom - scrollBarHeight, 0);
    },

    setHorizontalScrollbarWidth() {
      this.$refs.horizontalScrollbar?.setContentWidth?.(this.stickyTableWrapper?.scrollWidth ?? 0);
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
      this.stickyTableHeader?.classList.add('head-shadow');
    },

    removeShadowFromHeader() {
      this.stickyTableHeader?.classList.remove('head-shadow');
    },

    stickyScrollHandler() {
      this.clearFinishTimer();

      this.setHeaderWidth();
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

    scheduleStickyScrollHandler() {
      this.$nextTick(() => {
        window.requestAnimationFrame(() => {
          this.connectHorizontalScroll();
          this.stickyScrollHandler();
        });
      });
    },

    resetHeaderPosition() {
      this.translateY = 0;
      this.actionsTranslateY = 0;
      this.previousTranslateY = 0;
      this.previousScrollTranslateY = 0;

      if (this.stickyTableHeader) {
        this.stickyTableHeader.style.opacity = '';
        this.stickyTableHeader.style.transform = '';
        this.stickyTableHeader.style.width = '';
      }

      const stickyActions = this.getStickyActionsEl();

      if (stickyActions) {
        stickyActions.style.opacity = '';
        stickyActions.style.transform = '';
      }

      this.setHeaderPosition();
      this.clearFinishTimer();
      this.removeShadowFromHeader();
    },
  },
};
