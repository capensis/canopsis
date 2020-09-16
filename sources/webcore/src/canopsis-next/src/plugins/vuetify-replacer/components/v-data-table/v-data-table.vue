<script>
import { VDataTable } from 'vuetify/es5/components/VDataTable';

export default {
  extends: VDataTable,
  props: {
    isDisabledItem: {
      type: Function,
      default: item => !item,
    },
  },
  computed: {
    everyItem() {
      return this.filteredItems.length &&
        this.filteredItems.every(item => this.isDisabledItem(item) || this.isSelected(item));
    },
  },
};
</script>
