<script>
import { VSimpleTable } from 'vuetify/lib/components/VDataTable';

export default {
  extends: VSimpleTable,
  props: {
    ultraDense: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    classes() {
      return {
        'v-data-table--dense': this.dense,
        'v-data-table--ultra-dense': this.ultraDense,
        'v-data-table--fixed-height': !!this.height && !this.fixedHeader,
        'v-data-table--fixed-header': this.fixedHeader,
        'v-data-table--has-top': !!this.$slots.top,
        'v-data-table--has-bottom': !!this.$slots.bottom,
        ...this.themeClasses,
      };
    },
  },
};
</script>
