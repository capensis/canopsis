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
      scrollBarTranslateY: 0,
      stickyScrollTarget: null,
      translateY: 0,
    };
  },

  computed: {
    topBarHeight() {
      if (this.isFullscreenMode || this.isKioskMode) {
        return 0;
      }

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
    activeViewScreenMode() {
      if (this.someSticky) {
        this.scheduleStickyScrollHandler();
      }
    },

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

      window.addEventListener('scroll', this.stickyScrollEventHandler);
      window.addEventListener('resize', this.scheduleStickyScrollHandler);
      document.addEventListener('fullscreenchange', this.updateStickyScrollTarget);
      this.updateStickyScrollTarget();
      this.registeredScrollHandler = true;
    },

    unregisterScrollHandler() {
      window.removeEventListener('scroll', this.stickyScrollEventHandler);
      window.removeEventListener('resize', this.scheduleStickyScrollHandler);
      document.removeEventListener('fullscreenchange', this.updateStickyScrollTarget);
      this.stickyScrollTarget?.removeEventListener('scroll', this.stickyScrollEventHandler);
      this.stickyScrollTarget = null;
      this.registeredScrollHandler = false;
    },

    updateStickyScrollTarget() {
      const scrollTarget = this.$el?.closest('.view-fullscreen') ?? null;

      if (scrollTarget === this.stickyScrollTarget) {
        return;
      }

      this.stickyScrollTarget?.removeEventListener('scroll', this.stickyScrollEventHandler);
      scrollTarget?.addEventListener('scroll', this.stickyScrollEventHandler);
      this.stickyScrollTarget = scrollTarget;
      this.scheduleStickyScrollHandler();
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
      if (!this.stickyHeader || !this.stickyTableHeader) {
        return;
      }

      this.stickyTableHeader.style.opacity = '0';

      const stickyActions = this.getStickyActionsEl();

      if (stickyActions) {
        stickyActions.style.opacity = '0';
      }
    },

    startScrollingForScrollbar() {
      if (
        this.$refs.horizontalScrollbar &&
        this.previousScrollTranslateY !== this.scrollBarTranslateY
      ) {
        this.$refs.horizontalScrollbar.setOpacity('0');
      }
    },

    finishScrollingForHeader() {
      if (!this.stickyTableHeader || Number(this.stickyTableHeader.style.opacity)) {
        return;
      }

      this.stickyTableHeader.style.opacity = '1';

      const stickyActions = this.getStickyActionsEl();

      if (stickyActions) {
        stickyActions.style.opacity = '1';
      }
    },

    finishScrollingForScrollbar() {
      if (this.$refs.horizontalScrollbar?.$el?.style?.opacity === '0') {
        this.$refs.horizontalScrollbar.setOpacity('1');
      }
    },

    fullFinishScrolling() {
      this.finishScrollingForHeader();
      this.finishScrollingForScrollbar();
    },

    clearFinishTimer() {
      if (this.finishTimer) {
        clearTimeout(this.finishTimer);
      }
    },

    stickyScrollEventHandler() {
      this.clearFinishTimer();
      this.startScrollingForHeader();
      this.startScrollingForScrollbar();
      this.stickyScrollHandler();

      this.finishTimer = setTimeout(this.fullFinishScrolling, ALARMS_LIST_HEADER_OPACITY_DELAY);
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
      const { top: actionsTop = 0, height: actionsHeight = 0 } =
        this.getStickyActionsEl()?.getBoundingClientRect() ?? {};
      const offset = headerTop - this.translateY - this.topBarHeight - actionsHeight;
      const actionsOffset = actionsTop - this.actionsTranslateY - this.topBarHeight;

      this.translateY = Math.min(bodyHeight, Math.max(0, -offset));
      this.actionsTranslateY = Math.min(bodyHeight, Math.max(0, -actionsOffset));
    },

    calculateScrollPosition() {
      if (!this.stickyTableBody) {
        return;
      }

      const { bottom: bodyBottom } = this.stickyTableBody.getBoundingClientRect();
      const { height: scrollBarHeight = 0 } =
        this.$refs.horizontalScrollbar?.$el.getBoundingClientRect() ?? {};

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
      this.calculateStickyPositions();
      this.setHeaderPosition();

      if (this.stickyHorizontalScroll) {
        this.setHorizontalScrollbarWidth();
        this.setHorizontalScrollPosition();
      }

      if (!this.scrollBarTranslateY) {
        this.finishScrollingForScrollbar();
      }
    },

    scheduleStickyScrollHandler() {
      this.$nextTick(() => {
        window.requestAnimationFrame(() => {
          this.connectHorizontalScroll();
          this.setHeaderWidth();
          this.stickyScrollHandler();
        });
      });
    },

    resetHeaderPosition() {
      this.translateY = 0;
      this.actionsTranslateY = 0;
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
