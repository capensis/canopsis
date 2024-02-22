import { TOP_BAR_HEIGHT } from '@/config';
import { ALARMS_LIST_HEADER_OPACITY_DELAY } from '@/constants';

export const widgetHeaderStickyAlarmMixin = {
  props: {
    stickyHeader: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    topBarHeight() {
      return this.shownHeader ? TOP_BAR_HEIGHT : 0;
    },

    tableHeader() {
      return this.$el.querySelector('.alarms-list-table .v-data-table-header');
    },

    tableBody() {
      return this.$el.querySelector('.alarms-list-table table  tbody');
    },
  },

  watch: {
    stickyHeader(stickyHeader) {
      if (stickyHeader) {
        this.calculateHeaderOffsetPosition();
        this.setHeaderPosition();
        this.addShadowToHeader();

        window.addEventListener('scroll', this.changeHeaderPosition);
      } else {
        window.removeEventListener('scroll', this.changeHeaderPosition);

        this.resetHeaderPosition();
      }
    },
  },

  created() {
    this.actionsTranslateY = 0;
    this.translateY = 0;
    this.previousTranslateY = 0;
  },

  async mounted() {
    if (this.stickyHeader) {
      window.addEventListener('scroll', this.changeHeaderPosition);
    }
  },

  beforeDestroy() {
    window.removeEventListener('scroll', this.changeHeaderPosition);
  },

  methods: {
    startScrolling() {
      if (this.translateY !== this.previousTranslateY) {
        this.tableHeader.style.opacity = '0';

        if (this.$refs.actions) {
          this.$refs.actions.style.opacity = '0';
        }
      }

      this.scrooling = true;
    },

    finishScrolling() {
      if (!Number(this.tableHeader.style.opacity)) {
        this.tableHeader.style.opacity = '1.0';

        if (this.$refs.actions) {
          this.$refs.actions.style.opacity = '1.0';
        }
      }

      this.scrooling = false;
    },

    clearFinishTimer() {
      if (this.finishTimer) {
        clearTimeout(this.finishTimer);
      }
    },

    setHeaderPosition() {
      this.tableHeader.style.transform = `translateY(${this.translateY}px)`;

      if (this.$refs.actions) {
        this.$refs.actions.style.transform = `translateY(${this.actionsTranslateY}px)`;
      }
    },

    calculateHeaderOffsetPosition() {
      const { top: headerTop } = this.tableHeader.getBoundingClientRect();
      const { height: bodyHeight } = this.tableBody.getBoundingClientRect();
      const { top: actionsTop = 0, height: actionsHeight = 0 } = this.$refs.actions?.getBoundingClientRect() ?? {};

      const offset = headerTop - this.translateY - this.topBarHeight - actionsHeight;
      const actionsOffset = actionsTop - this.actionsTranslateY - this.topBarHeight;

      this.previousTranslateY = this.actionsTranslateY;
      this.translateY = Math.min(bodyHeight, Math.max(0, -offset));
      this.actionsTranslateY = Math.min(bodyHeight, Math.max(0, -actionsOffset));
    },

    addShadowToHeader() {
      this.tableHeader.classList.add('head-shadow');
    },

    removeShadowFromHeader() {
      this.tableHeader.classList.remove('head-shadow');
    },

    changeHeaderPosition() {
      this.clearFinishTimer();

      this.calculateHeaderOffsetPosition();
      this.setHeaderPosition();

      if (!this.actionsTranslateY || !this.translateY) {
        this.removeShadowFromHeader();
        this.finishScrolling();

        return;
      }

      if (!this.scrooling) {
        this.addShadowToHeader();
        this.startScrolling();
      }

      this.finishTimer = setTimeout(this.finishScrolling, ALARMS_LIST_HEADER_OPACITY_DELAY);
    },

    resetHeaderPosition() {
      this.translateY = 0;
      this.actionsTranslateY = 0;
      this.previousTranslateY = 0;

      this.setHeaderPosition();
      this.clearFinishTimer();
      this.removeShadowFromHeader();
    },
  },
};
