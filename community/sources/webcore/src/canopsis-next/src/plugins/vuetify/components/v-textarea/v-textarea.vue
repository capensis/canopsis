<script>
import VTextarea from 'vuetify/lib/components/VTextarea';

import { getStringLinesCount } from '@/helpers/string';

export default {
  extends: VTextarea,
  data() {
    return {
      $patchedTextareaResizeObserver: null,
      patchedTextareaRafId: 0,
    };
  },
  watch: {
    autoGrow: {
      immediate: false,
      handler(value) {
        if (!value) {
          this.patchedTextareaDestroyResizeObserver();
          return;
        }

        this.$nextTick(() => {
          this.patchedTextareaRecalculateHeight();
          this.patchedTextareaInitResizeObserver();
        });
      },
    },
  },
  mounted() {
    if (!this.autoGrow) {
      return;
    }

    this.$nextTick(() => {
      this.patchedTextareaRecalculateHeight();
      this.patchedTextareaInitResizeObserver();
    });
  },

  updated() {
    if (!this.autoGrow) {
      this.patchedTextareaDestroyResizeObserver();
      return;
    }

    this.$nextTick(() => {
      this.patchedTextareaRecalculateHeight();
      this.patchedTextareaInitResizeObserver();
    });
  },

  beforeDestroy() {
    this.patchedTextareaDestroyResizeObserver();

    if (this.patchedTextareaRafId) {
      window.cancelAnimationFrame(this.patchedTextareaRafId);
      this.patchedTextareaRafId = 0;
    }
  },
  methods: {
    calculateInputHeight() {
      const { input } = this.$refs;

      if (!input) return;
      input.style.height = '0';

      const computedStyle = window.getComputedStyle(input);
      const lineHeight = parseFloat(computedStyle.lineHeight);
      const linesCount = getStringLinesCount(this.value);
      const height = Math.ceil(lineHeight * linesCount);
      const minHeight = parseInt(this.rows, 10) * parseFloat(this.rowHeight);

      input.style.height = `${Math.max(height, minHeight)}px`;
    },

    onInput(event) {
      VTextarea.options.methods.onInput.call(this, event);

      if (this.autoGrow) {
        this.patchedTextareaScheduleRecalculate();
      }
    },

    patchedTextareaGetNativeTextarea() {
      return this.$el?.querySelector?.('textarea') ?? null;
    },

    patchedTextareaScheduleRecalculate() {
      if (!this.autoGrow) {
        return;
      }

      if (this.patchedTextareaRafId) {
        window.cancelAnimationFrame(this.patchedTextareaRafId);
      }

      this.patchedTextareaRafId = window.requestAnimationFrame(() => {
        this.patchedTextareaRafId = 0;
        this.patchedTextareaRecalculateHeight();
      });
    },

    patchedTextareaRecalculateHeight() {
      if (!this.autoGrow) {
        return;
      }

      const textarea = this.patchedTextareaGetNativeTextarea();

      if (!textarea) {
        return;
      }

      const computedStyle = window.getComputedStyle(textarea);
      const borderTopWidth = parseFloat(computedStyle.borderTopWidth || '0') || 0;
      const borderBottomWidth = parseFloat(computedStyle.borderBottomWidth || '0') || 0;
      const verticalBorderWidth = borderTopWidth + borderBottomWidth;

      textarea.style.height = 'auto';
      textarea.style.overflowY = 'auto';

      let nextHeight = textarea.scrollHeight + verticalBorderWidth;

      const rows = Number(this.rows || 0);

      if (rows > 0) {
        const lineHeight = parseFloat(computedStyle.lineHeight || '0') || 24;
        const paddingTop = parseFloat(computedStyle.paddingTop || '0') || 0;
        const paddingBottom = parseFloat(computedStyle.paddingBottom || '0') || 0;
        const minHeight = rows * lineHeight + paddingTop + paddingBottom + verticalBorderWidth;

        nextHeight = Math.max(nextHeight, minHeight);
      }

      const maxRows = Number(this.maxRows || 0);

      if (maxRows > 0) {
        const lineHeight = parseFloat(computedStyle.lineHeight || '0') || 24;
        const paddingTop = parseFloat(computedStyle.paddingTop || '0') || 0;
        const paddingBottom = parseFloat(computedStyle.paddingBottom || '0') || 0;
        const maxHeight = maxRows * lineHeight + paddingTop + paddingBottom + verticalBorderWidth;

        nextHeight = Math.min(nextHeight, maxHeight);
        textarea.style.overflowY = textarea.scrollHeight > maxHeight ? 'auto' : 'hidden';
      }

      textarea.style.height = `${nextHeight}px`;
    },

    patchedTextareaInitResizeObserver() {
      if (!this.autoGrow || !ResizeObserver || !this.$el || this.$patchedTextareaResizeObserver) {
        return;
      }

      // eslint-disable-next-line no-underscore-dangle
      this.$el.__patchedTextareaPreviousWidth = this.$el.clientWidth;

      this.$patchedTextareaResizeObserver = new ResizeObserver(() => {
        const currentWidth = this.$el.clientWidth;
        // eslint-disable-next-line no-underscore-dangle
        const previousWidth = this.$el.__patchedTextareaPreviousWidth;

        if (previousWidth !== currentWidth) {
          // eslint-disable-next-line no-underscore-dangle
          this.$el.__patchedTextareaPreviousWidth = currentWidth;
          this.patchedTextareaScheduleRecalculate();
        }
      });

      this.$patchedTextareaResizeObserver.observe(this.$el);
    },

    patchedTextareaDestroyResizeObserver() {
      if (!this.$patchedTextareaResizeObserver) {
        return;
      }

      this.$patchedTextareaResizeObserver.disconnect();
      this.$patchedTextareaResizeObserver = null;
    },
  },
};
</script>
