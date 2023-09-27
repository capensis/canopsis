<script>
import VChip from 'vuetify/es5/components/VChip';
import VIcon from 'vuetify/es5/components/VIcon';

export default {
  extends: VChip,
  props: {
    closeIcon: {
      type: String,
      default: '$vuetify.icons.delete',
    },
  },
  methods: {
    genClose(h) {
      const data = {
        staticClass: 'v-chip__close',
        on: {
          click: (e) => {
            e.stopPropagation();
            this.$emit('input', false);
          },
        },
      };

      return h('div', data, [
        h(VIcon, this.closeIcon || '$vuetify.icons.delete'),
      ]);
    },
  },
};
</script>
