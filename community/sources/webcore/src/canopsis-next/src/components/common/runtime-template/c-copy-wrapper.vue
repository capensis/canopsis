<script>
export default {
  props: {
    value: {
      type: String,
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
  },
  render(h) {
    const slots = this.$slots.default || [];
    const [slot] = slots;

    let tag = 'span';
    let children = slots;
    let attributes = {
      directives: [{
        name: 'clipboard',
        rawName: 'v-clipboard:copy',
        value: this.value,
        arg: 'copy',
      },
      {
        name: 'clipboard',
        rawName: 'v-clipboard:success',
        value: this.showSuccessPopup,
        arg: 'success',
      },
      {
        name: 'clipboard',
        rawName: 'v-clipboard:error',
        value: this.showErrorPopup,
        arg: 'error',
      }],
    };

    if (slots.length === 1 && slot?.tag) {
      tag = slot?.tag;
      attributes = { ...attributes, ...slot.data };
      children = slot.children;
    }

    return h(tag, attributes, children);
  },
};
</script>
