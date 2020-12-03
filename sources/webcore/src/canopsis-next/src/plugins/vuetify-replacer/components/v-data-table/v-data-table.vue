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
    activeItems() {
      return this.filteredItems.filter(item => !this.isDisabledItem(item));
    },

    everyItem() {
      return this.activeItems.length && this.activeItems.every(this.isSelected);
    },
  },
};
</script>
