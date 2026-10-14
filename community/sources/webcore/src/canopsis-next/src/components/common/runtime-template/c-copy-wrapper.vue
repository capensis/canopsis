<script>
import { cloneVNode, h } from 'vue';

import { writeTextToClipboard } from '@/helpers/clipboard';

export default {
  props: {
    value: {
      type: [String, Number],
      default: '',
    },
  },

  methods: {
    showSuccessPopup() {
      this.$popups.success({ text: this.$t('popups.copySuccess') });
    },

    showErrorPopup() {
      this.$popups.error({ text: this.$t('popups.copyError') });
    },

    async copy() {
      try {
        await writeTextToClipboard(String(this.value));
        this.showSuccessPopup();
      } catch {
        this.showErrorPopup();
      }
    },

    prepareVNode(vnode) {
      return cloneVNode(vnode, {
        onClick: (event) => {
          vnode.props?.onClick?.(event);
          this.copy();
        },
      });
    },
  },

  render() {
    const slots = this.$slots.default?.() ?? [];

    if (slots.length === 1) {
      return this.prepareVNode(slots[0]);
    }

    return h(
      'span',
      {
        onClick: this.copy,
      },
      slots,
    );
  },
};
</script>
